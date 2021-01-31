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

package v1alpha1

import (
	"fmt"
	"net"

	"github.com/95ulisse/dns-operator/pkg/dnsname"
	"github.com/miekg/dns"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSRecordSpec defines the desired state of DNSRecord
type DNSRecordSpec struct {
	// Reference to the DNSProvider managing this DNSRecord.
	ProviderRef ObjectReference `json:"providerRef"`

	// Name of the DNS record.
	// This field is required.
	Name dnsname.Name `json:"name"`

	// Content of the DNS record. The meaning of the content field depends on the type of record.
	// This field is required.
	Content DNSRecordContent `json:"content"`

	// TTL in seconds of the DNS record. Defaults to 1h.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=604800
	// +optional
	TTLSeconds *uint32 `json:"ttlSeconds,omitempty"`

	// Specifies how to treat deletion of this DNSRecord.
	// Valid values are:
	// - "Delete" (default): actually delete the corresponding DNS record managed by this resource;
	// - "Retain": do not delete the actual DNS record managed by this resource.
	// +optional
	DeletionPolicy *DeletionPolicy `json:"deletionPolicy,omitempty"`
}

// DNSRecordContent represents the actual contents of a DNS record.
// Only one of these can be set.
type DNSRecordContent struct {
	// A record.
	// +kubebuilder:validation:Format=ipv4
	// +optional
	A *string `json:"a,omitempty"`

	// AAAA record.
	// +kubebuilder:validation:Format=ipv6
	// +optional
	AAAA *string `json:"aaaa,omitempty"`

	// CNAME record.
	// +kubebuilder:validation:MinLength=0
	// +optional
	CNAME *string `json:"cname,omitempty"`

	// TXT record.
	// +optional
	TXT *string `json:"txt,omitempty"`

	// NS record.
	// +kubebuilder:validation:MinLength=0
	// +optional
	NS *string `json:"ns,omitempty"`

	// MX record.
	// +optional
	MX *MXRecordContent `json:"mx,omitempty"`
}

// MXRecordContent represents the contents of an MX DNS record.
type MXRecordContent struct {
	// Name pointed by the MX record.
	// +kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	// Priority of the MX record.
	// +kubebuilder:validation:Minimum=0
	Priority int64 `json:"priority"`
}

// DNSRecordStatus defines the observed state of DNSRecord
type DNSRecordStatus struct {
	StatusWithConditions `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="RR Name",type="string",JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="RR Data",type="string",JSONPath=`.spec.content`

// DNSRecord is the Schema for the dnsrecords API
type DNSRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSRecordSpec   `json:"spec,omitempty"`
	Status DNSRecordStatus `json:"status,omitempty"`
}

// ToRR builds a new dns.RR equivalent to the Spec of this DNSRecord resource.
func (record *DNSRecord) ToRR() (dns.RR, error) {

	// Prepare a common header
	header := dns.RR_Header{
		Name:  record.Spec.Name.ToFQDN().String(),
		Class: dns.ClassINET,
		Ttl:   3600,
	}
	if record.Spec.TTLSeconds != nil {
		header.Ttl = *record.Spec.TTLSeconds
	}

	// A record
	if record.Spec.Content.A != nil {
		ip := net.ParseIP(*record.Spec.Content.A)
		if ip == nil {
			return nil, fmt.Errorf("Invalid IPv4 address %s", *record.Spec.Content.A)
		}
		ip = ip.To4()
		if ip == nil {
			return nil, fmt.Errorf("Invalid IPv4 address %s", *record.Spec.Content.A)
		}

		rr := new(dns.A)
		rr.Hdr = header
		rr.Hdr.Rrtype = dns.TypeA
		rr.A = ip
		return rr, nil
	}

	return nil, fmt.Errorf("Unsupported DNS record")
}

// +kubebuilder:object:root=true

// DNSRecordList contains a list of DNSRecord
type DNSRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSRecord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSRecord{}, &DNSRecordList{})
}
