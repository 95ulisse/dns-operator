package providers

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// Dummy DNS provider for debugging.
type Dummy struct {
	log logr.Logger
}

// UpdateRecord dummy noop.
func (dummy *Dummy) UpdateRecord(record *dnsv1alpha1.DNSRecord) error {
	dummy.log.Info("Update successful")
	return nil
}

// DeleteRecord dummy noop.
func (dummy *Dummy) DeleteRecord(record *dnsv1alpha1.DNSRecord) error {
	dummy.log.Info("Delete successful")
	return nil
}

func init() {
	RegisterProviderFactory(func(resource *dnsv1alpha1.DNSProvider, _ client.Client, log logr.Logger) (Provider, error) {

		// Skip this provider if not configured
		if resource.Spec.Dummy == nil {
			return nil, nil
		}

		dummy := Dummy{
			log: log.
				WithName("providers").WithName("Dummy").
				WithValues("dnsprovider", fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)),
		}
		return &dummy, nil

	})
}
