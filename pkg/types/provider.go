package types

import (
	"github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/dnsname"
)

// Provider is a generic DNS provider which knows how to talk to a backend to reconcile DNS records.
type Provider interface {
	Zones() []dnsname.Name
	UpdateRecord(zone dnsname.Name, rrset v1alpha1.DNSRecord) error
	DeleteRecord(zone dnsname.Name, rrset v1alpha1.DNSRecord) error
}
