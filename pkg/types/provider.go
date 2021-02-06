package types

import (
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	"github.com/miekg/dns"
)

// Provider is a generic DNS provider which knows how to talk to a backend to reconcile DNS records.
type Provider interface {
	Zones() []dnsname.Name
	UpdateRecord(zone dnsname.Name, rrset []dns.RR) error
	DeleteRecord(zone dnsname.Name, rrset []dns.RR) error
}
