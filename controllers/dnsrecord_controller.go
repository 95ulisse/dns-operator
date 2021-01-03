/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/api/v1alpha1"
)

// DNSRecordReconciler reconciles a DNSRecord object
type DNSRecordReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsproviders,verbs=get;list;watch
// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsrecords,verbs=get;list;watch
// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsrecords/status,verbs=get;update;patch

// Reconcile performs an iteration of the reconcile loop for a DNSRecord.
func (r *DNSRecordReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("dnsrecord", req.NamespacedName)

	log.V(1).Info("Starting reconcile loop")
	defer log.V(1).Info("Finish reconcile loop")

	// Retrieve the record by name
	var record dnsv1alpha1.DNSRecord
	if err := r.Get(ctx, req.NamespacedName, &record); err != nil {
		log.Error(err, "Unable to fetch DNSRecord")
		// We'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Retrieve the provider
	refName := record.Spec.ProviderRef.Name
	refNamespace := record.Spec.ProviderRef.Namespace
	if refNamespace == nil {
		refNamespace = &record.Namespace
	}
	var providerKey types.NamespacedName = types.NamespacedName{
		Name:      refName,
		Namespace: *refNamespace,
	}
	var provider dnsv1alpha1.DNSProvider
	if err := r.Get(ctx, providerKey, &provider); err != nil {
		log.Error(err, "Cannot find DNSProvider", "dnsprovider", providerKey)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.V(1).Info(fmt.Sprintf("Found provider %v", providerKey))

	return ctrl.Result{}, nil
}

// listRecordsUsingProvider returns a list of the names of DNSRecords resources that reference the given DNSProvider.
func (r *DNSRecordReconciler) listRecordsUsingProvider(provider handler.MapObject) []ctrl.Request {

	// Filter all the DNSRecords using the provider referencing the given provider by name.
	// We do not check the namespace of the records because records can leave the namespace
	// of the reference unspecified, and we cannot filter on un set fields.
	//
	// TODO: This can be solved using a defaulting webhook that automatically applies
	// the correct DNSProvider ref namespace, so that we can index it.
	listOptions := []client.ListOption{
		client.MatchingField(".spec.providerRef.name", provider.Meta.GetName()),
	}
	var list dnsv1alpha1.DNSRecordList
	if err := r.List(context.Background(), &list, listOptions...); err != nil {
		r.Log.Error(
			err,
			"Cannot list DNSRecords impacted by a change to DNSProvider",
			"dnsprovider", fmt.Sprintf("%s/%s", provider.Meta.GetNamespace(), provider.Meta.GetName()),
		)
		return nil
	}

	var res []ctrl.Request
	for _, record := range list.Items {

		// By default, if the provider reference of this record does not include a namespace,
		// use the same namespace of the record itself.
		refName := record.Spec.ProviderRef.Name
		refNamespace := record.Spec.ProviderRef.Namespace
		if refNamespace == nil {
			refNamespace = &record.Namespace
		}

		// Select this record for reconciling if the provider ref matches the changed provider
		if refName == provider.Meta.GetName() && *refNamespace == provider.Meta.GetNamespace() {
			res = append(res, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      record.Name,
					Namespace: record.Namespace,
				},
			})
		}

	}

	r.Log.V(1).Info(
		"Enqueued reconciling of DNSRecord due to a change to DNSProvider",
		"count", len(res),
		"dnsprovider", fmt.Sprintf("%s/%s", provider.Meta.GetNamespace(), provider.Meta.GetName()),
	)

	return res
}

// SetupWithManager registers the DNSRecord controller with the given Manager.
func (r *DNSRecordReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// Index DNSRecords by the name of the provider they used
	mgr.GetFieldIndexer().IndexField(
		&dnsv1alpha1.DNSRecord{},
		".spec.providerRef.name",
		func(obj runtime.Object) []string {
			providerName := obj.(*dnsv1alpha1.DNSRecord).Spec.ProviderRef.Name
			if providerName == "" {
				return nil
			}
			return []string{providerName}
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1alpha1.DNSRecord{}).
		Watches(
			&source.Kind{Type: &dnsv1alpha1.DNSProvider{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.listRecordsUsingProvider),
			},
		).
		Complete(r)
}
