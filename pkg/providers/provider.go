package providers

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// Provider is a generic DNS provider which knows how to talk to a backend to reconcile DNS records.
type Provider interface {
	UpdateRecord(record *dnsv1alpha1.DNSRecord) error
	DeleteRecord(record *dnsv1alpha1.DNSRecord) error
}

var factories []func(*dnsv1alpha1.DNSProvider, client.Client, logr.Logger) (Provider, error)

// FromKubernetesResource builds a new Provider from the configuration in the kubernetes DNSProvider resource.
func FromKubernetesResource(provider *dnsv1alpha1.DNSProvider, client client.Client, log logr.Logger) (Provider, error) {
	for _, f := range factories {
		res, err := f(provider, client, log)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, fmt.Errorf("Unable to parse configuration of DNSProvider %s/%s", provider.Namespace, provider.Name)
}

// RegisterProviderFactory registers a new factory function used to build an actual Provider
// from the description in the kubernetes DNSProvider resource.
// The factory function returns (nil, nil) if it does not know how to build a provider.
func RegisterProviderFactory(factory func(*dnsv1alpha1.DNSProvider, client.Client, logr.Logger) (Provider, error)) {
	factories = append(factories, factory)
}
