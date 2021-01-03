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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSRecordSpec defines the desired state of DNSRecord
type DNSRecordSpec struct {
	// Reference to the DNSProvider managing this DNSRecord.
	ProviderRef ObjectReference `json:"providerRef"`

	// Name of the DNS record.
	// This field is required.
	// +kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	// Content of the DNS record. The meaning of the content field depends on the type of record.
	// This field is required.
	Content DNSRecordContent `json:"content"`

	// TTL in seconds of the DNS record. Defaults to 1h.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=604800
	// +optional
	TTLSeconds *int64 `json:"ttlSeconds,omitempty"`

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
