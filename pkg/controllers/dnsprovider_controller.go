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
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/providers"
	"github.com/95ulisse/dns-operator/pkg/types"
)

// DNSProviderReconciler reconciles a DNSProvider object
type DNSProviderReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Context *types.ControllerContext
}

// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsproviders,verbs=get;list;watch
// +kubebuilder:rbac:groups=dns.k8s.marcocameriero.net,resources=dnsproviders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile performs an iteration of the reconcile loop for a DNSProvider.
func (r *DNSProviderReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := r.Context.RootContext
	log := r.Log.WithValues("dnsprovider", req.NamespacedName)

	log.V(1).Info("Starting reconcile loop")
	defer log.V(1).Info("Finish reconcile loop")

	// Retrieve the provider by name
	var resource dnsv1alpha1.DNSProvider
	if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
		// We get not-found erros after object deletion.
		// Remove the provider from the global context in that case.
		if apierrors.IsNotFound(err) {
			r.Context.RemoveProvider(req.NamespacedName.String())
			log.V(1).Info("Removed provider")
		} else {
			log.Error(err, "Unable to fetch DNSProvider")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Mark the provider as non ready
	resource.Status.SetCondition(&dnsv1alpha1.Condition{
		Type:    dnsv1alpha1.ReadyCondition,
		Status:  dnsv1alpha1.FalseStatus,
		Reason:  "Configuring",
		Message: "Configuring the provider",
	})
	if err := r.Status().Update(ctx, &resource); err != nil {
		log.Error(err, "Cannot update resource status")
		return ctrl.Result{}, err
	}

	// Build the actual provider and store it in the shared context
	provider, err := providers.ProviderFor(r.Context, &resource)
	if err != nil {
		log.Error(err, "Cannot build provider")

		// Mark the provider as non ready
		resource.Status.SetCondition(&dnsv1alpha1.Condition{
			Type:    dnsv1alpha1.ReadyCondition,
			Status:  dnsv1alpha1.FalseStatus,
			Reason:  "Error",
			Message: err.Error(),
		})
		if err := r.Status().Update(ctx, &resource); err != nil {
			log.Error(err, "Cannot update resource status")
			return ctrl.Result{}, err
		}

		// Record an event
		r.Context.EventRecorder.Event(&resource, "Warning", "Error", err.Error())

		return ctrl.Result{}, err
	}
	r.Context.SetProvider(req.NamespacedName.String(), provider)
	log.Info("Provider updated")

	// Mark the provider as ready
	resource.Status.SetCondition(&dnsv1alpha1.Condition{
		Type:    dnsv1alpha1.ReadyCondition,
		Status:  dnsv1alpha1.TrueStatus,
		Reason:  "Ready",
		Message: "Ready to register DNS records",
	})
	if err := r.Status().Update(ctx, &resource); err != nil {
		log.Error(err, "Cannot update resource status")
		return ctrl.Result{}, err
	}

	// Record an event
	r.Context.EventRecorder.Event(&resource, "Normal", "Ready", "Ready to register DNS records")

	return ctrl.Result{}, nil
}

// SetupWithManager registers the DNSProvider controller with the given Manager.
func (r *DNSProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1alpha1.DNSProvider{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
