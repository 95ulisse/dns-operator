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

// Reconcile performs an iteration of the reconcile loop for a DNSProvider.
func (r *DNSProviderReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := r.Context.RootContext
	log := r.Log.WithValues("dnsprovider", req.NamespacedName)

	log.V(1).Info("Starting reconcile loop")
	defer log.V(1).Info("Finish reconcile loop")

	// Remove the provider from the global context
	r.Context.RemoveProvider(req.NamespacedName.String())
	log.V(1).Info("Removed provider")

	// Retrieve the provider by name
	var resource dnsv1alpha1.DNSProvider
	if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
		// We'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		if !apierrors.IsNotFound(err) {
			log.Error(err, "Unable to fetch DNSProvider")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Build the actual provider and store it in the shared context
	provider, err := providers.ProviderFor(r.Context, &resource)
	if err != nil {
		log.Error(err, "Cannot build provider")
		return ctrl.Result{}, err
	}
	r.Context.SetProvider(req.NamespacedName.String(), provider)
	log.Info("Provider updated")

	return ctrl.Result{}, nil
}

// SetupWithManager registers the DNSProvider controller with the given Manager.
func (r *DNSProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1alpha1.DNSProvider{}).
		Complete(r)
}
