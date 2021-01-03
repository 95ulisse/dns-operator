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

	log.V(1).Info(fmt.Sprintf("Starting reconcile loop for %v", req.NamespacedName))
	defer log.V(1).Info(fmt.Sprintf("Finish reconcile loop for %v", req.NamespacedName))

	// Retrieve the record by name
	var record dnsv1alpha1.DNSRecord
	if err := r.Get(ctx, req.NamespacedName, &record); err != nil {
		log.Error(err, "Unable to fetch DNSRecord")
		// We'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// listRecordsUsingProvider returns a list of the names of DNSRecords resources that reference the given DNSProvider.
func (r *DNSRecordReconciler) listRecordsUsingProvider(obj handler.MapObject) []ctrl.Request {
	listOptions := []client.ListOption{
		// matching our index
		client.MatchingField(".spec.providerRef.name", obj.Meta.GetName()),
		client.MatchingField(".spec.providerRef.namespace", obj.Meta.GetNamespace()),
	}
	var list dnsv1alpha1.DNSRecordList
	if err := r.List(context.Background(), &list, listOptions...); err != nil {
		r.Log.Error(err, "Cannot list DNSRecords impacted by a change to the DNSProvider %s/%s", obj.Meta.GetNamespace(), obj.Meta.GetName())
		return nil
	}
	res := make([]ctrl.Request, len(list.Items))
	for i, record := range list.Items {
		res[i].Name = record.Name
		res[i].Namespace = record.Namespace
	}
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
	mgr.GetFieldIndexer().IndexField(
		&dnsv1alpha1.DNSRecord{},
		".spec.providerRef.namespace",
		func(obj runtime.Object) []string {
			providerName := obj.(*dnsv1alpha1.DNSRecord).Spec.ProviderRef.Namespace
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
