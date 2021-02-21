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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	helpers "github.com/95ulisse/dns-operator/pkg/helpers"
	"github.com/95ulisse/dns-operator/pkg/types"
)

const finalizerName = "dns.k8s.marcocameriero.net/finalizer"

// DNSRecordReconciler reconciles a DNSRecord object
type DNSRecordReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Context *types.ControllerContext
}

// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsrecords,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsrecords/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch;update

// Reconcile performs an iteration of the reconcile loop for a DNSRecord.
func (r *DNSRecordReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := r.Context.RootContext
	log := r.Log.WithValues("dnsrecord", req.NamespacedName)

	log.V(1).Info("Starting reconcile loop")
	defer log.V(1).Info("Finish reconcile loop")

	// Step 1: Retrive the DNSRecord to reconcile
	// ==========================================

	// Retrieve the record by name
	var record dnsv1alpha1.DNSRecord
	if err := r.Get(ctx, req.NamespacedName, &record); err != nil {
		// We'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		if !apierrors.IsNotFound(err) {
			log.Error(err, "Unable to fetch DNSRecord")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Mark the record as not ready
	record.Status.SetCondition(&dnsv1alpha1.Condition{
		Type:    dnsv1alpha1.ReadyCondition,
		Status:  dnsv1alpha1.FalseStatus,
		Reason:  "NotReady",
		Message: "",
	})
	if err := r.Status().Update(ctx, &record); err != nil {
		log.Error(err, "Cannot update resource status")
		return ctrl.Result{}, err
	}

	// Step 2: Retrieve the referenced DNSProvider
	// ===========================================

	// Retrieve the provider
	refName := record.Spec.ProviderRef.Name
	refNamespace := record.Spec.ProviderRef.Namespace
	if refNamespace == nil {
		refNamespace = &record.Namespace
	}
	providerNamespacedName := fmt.Sprintf("%s/%s", *refNamespace, refName)
	var provider types.Provider
	providerFound := r.Context.GetProvider(providerNamespacedName, &provider)

	// Check that the provider manages a zone containing this record
	var zone dnsname.Name
	if providerFound {
		if !getMatchingZone(provider.Zones(), record.Spec.Name, &zone) {
			err := fmt.Errorf("Provider %s does not support a zone matching record %s", providerNamespacedName, record.Spec.Name.String())
			return ctrl.Result{}, err
		}
	}

	// Step 3: Process any pending finalizers
	// ======================================

	// Examine DeletionTimestamp to determine if object is under deletion
	if record.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// lets add the finalizer and update the object. This is equivalent to registering our finalizer.
		if !helpers.ContainsString(record.ObjectMeta.Finalizers, finalizerName) {
			record.ObjectMeta.Finalizers = append(record.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, &record); err != nil {
				return ctrl.Result{}, err
			}
			log.V(1).Info("Finalizer registered")
		}
	} else {

		// The object is being deleted, so execute the finalizer
		if helpers.ContainsString(record.ObjectMeta.Finalizers, finalizerName) {

			// Actually delete the record from the provider only if the user does not want us to retain the actual record
			if record.Spec.DeletionPolicy == nil || *record.Spec.DeletionPolicy == dnsv1alpha1.DeletePolicy {

				if !providerFound {
					err := fmt.Errorf("Cannot find DNSProvider %s", providerNamespacedName)
					log.Error(err, "Cannot delete DNSRecord")
					return ctrl.Result{}, err
				}

				log.V(1).Info("Deleting record")

				if err := provider.DeleteRecord(zone, record); err != nil {
					log.Error(err, "Cannot delete DNSRecord")
					return ctrl.Result{}, err
				}

			}

			// remove our finalizer from the list and update it.
			record.ObjectMeta.Finalizers = helpers.RemoveString(record.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, &record); err != nil {
				return ctrl.Result{}, err
			}

			log.Info("Record deleted")
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Step 4: Update the DNS record
	// =============================

	if !providerFound {
		err := fmt.Errorf("Cannot find DNSProvider %s", providerNamespacedName)
		log.Error(err, "Cannot update DNSRecord")
		return ctrl.Result{}, err
	}

	// Let the magic happen
	if err := provider.UpdateRecord(zone, record); err != nil {
		log.Error(err, "Cannot update update DNS record")
		return ctrl.Result{}, err
	}

	// Mark the record as ready
	record.Status.SetCondition(&dnsv1alpha1.Condition{
		Type:    dnsv1alpha1.ReadyCondition,
		Status:  dnsv1alpha1.TrueStatus,
		Reason:  "Ready",
		Message: "DNS record registered",
	})
	if err := r.Status().Update(ctx, &record); err != nil {
		log.Error(err, "Cannot update resource status")
		return ctrl.Result{}, err
	}

	log.Info("Successfully updated record")

	r.Context.EventRecorder.Event(&record, "Normal", "Registered", "DNS record correclty registered")

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
				NamespacedName: k8stypes.NamespacedName{
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

func getMatchingZone(zones []dnsname.Name, record dnsname.Name, out *dnsname.Name) bool {

	// Filter only the zones containing the target record
	found := make([]dnsname.Name, 0, 1)
	for _, z := range zones {
		if record.IsChildOf(&z) {
			found = append(found, z)
		}
	}

	if len(found) == 0 {
		return false
	}

	// Find the longest name
	longest := found[0]
	for i, z := range found {
		if i > 0 && len(longest.ToFQDN().String()) < len(z.ToFQDN().String()) {
			longest = z
		}
	}

	*out = longest
	return true

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
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
