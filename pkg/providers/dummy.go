package providers

import (
	"github.com/go-logr/logr"
	"github.com/miekg/dns"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	"github.com/95ulisse/dns-operator/pkg/types"
)

// Dummy DNS provider for debugging.
type Dummy struct {
	log   logr.Logger
	zones []dnsname.Name
}

// NewDummy creates a new instance of the Dummy provider.
func NewDummy(log logr.Logger, zones []dnsname.Name) *Dummy {
	return &Dummy{
		log:   log.WithName("providers").WithName("Dummy"),
		zones: zones,
	}
}

// Zones returns a slice containing the DNS zones managed by this provider.
func (dummy *Dummy) Zones() []dnsname.Name {
	return dummy.zones
}

// UpdateRecord dummy noop.
func (dummy *Dummy) UpdateRecord(zone dnsname.Name, record dns.RR) error {
	dummy.log.Info("Updating ")
	return nil
}

// DeleteRecord dummy noop.
func (dummy *Dummy) DeleteRecord(zone dnsname.Name, record dns.RR) error {
	dummy.log.Info("Delete successful")
	return nil
}

func init() {
	RegisterProviderConstructor("dummy", func(ctx *types.ControllerContext, resource *dnsv1alpha1.DNSProvider) (types.Provider, error) {
		return NewDummy(ctx.Log, resource.Spec.Zones), nil
	})
}
