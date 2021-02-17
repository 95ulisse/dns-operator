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

	"github.com/95ulisse/dns-operator/pkg/dnsname"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSProviderSpec defines the desired state of DNSProvider.
// Only one of the providers can be configured.
type DNSProviderSpec struct {

	// DNS zones handled by this provider.
	// At least one zone must be present.
	// +kubebuilder:validation:MinItems=1
	Zones []dnsname.Name `json:"zones"`

	// Dummy provider used for debugging.
	// +optional
	Dummy *bool `json:"dummy,omitempty"`

	// Use RFC2136 ("Dynamic Updates in the Domain Name System") (https://datatracker.ietf.org/doc/rfc2136/) to manage records.
	// +optional
	RFC2136 *DNSProviderRFC2136 `json:"rfc2136,omitempty"`

	// Use Cloudflare to manage records.
	// +optional
	Cloudflare *DNSProviderCloudflare `json:"cloudflare,omitempty"`
}

// DNSProviderRFC2136 is a structure containing the configuration for RFC2136 DNS provider.
type DNSProviderRFC2136 struct {
	// The IP address or hostname of an authoritative DNS server supporting
	// RFC2136 in the form host:port. If the host is an IPv6 address it must be
	// enclosed in square brackets (e.g [2001:db8::1]) ; port is optional.
	// This field is required.
	Nameserver string `json:"nameserver"`

	// The name of the secret containing the TSIG value.
	// If ``tsigKeyName`` is defined, this field is required.
	// +optional
	TSIGSecretRef *SecretReference `json:"tsigSecretRef,omitempty"`

	// The TSIG Key name configured in the DNS.
	// If ``tsigSecretSecretRef`` is defined, this field is required.
	// +optional
	TSIGKeyName *string `json:"tsigKeyName,omitempty"`

	// The TSIG Algorithm configured in the DNS supporting RFC2136. Used only
	// when ``tsigSecretSecretRef`` and ``tsigKeyName`` are defined.
	// Supported values are (case-insensitive): ``HMACMD5``,
	// ``HMACSHA1``, ``HMACSHA256`` or ``HMACSHA512``.
	// +optional
	TSIGAlgorithm *string `json:"tsigAlgorithm,omitempty"`
}

// DNSProviderCloudflare is a structure containing the configuration of the Cloudflare provider.
// The Cloudflare provider can be configured using either an API Token, or an API Key.
type DNSProviderCloudflare struct {
	// Email owner of the Cloudflare account, required only if using an API Key.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Format=email
	// +optional
	Email *string `json:"email,omitempty"`

	// Reference to a secret containing the API Token to use for authentication.
	// One between `apiTokenSecretRef` and `apiKeySecretRef` must be present.
	// +optional
	APITokenSecretRef *SecretReference `json:"apiTokenSecretRef,omitempty"`

	// Reference to a secret containing the API Key to use for authentication.
	// One between `apiTokenSecretRef` and `apiKeySecretRef` must be present.
	// +optional
	APIKeySecretRef *SecretReference `json:"apiKeySecretRef,omitempty"`
}

// DNSProviderStatus defines the observed state of DNSProvider
type DNSProviderStatus struct {
	StatusWithConditions `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DNSProvider is the Schema for the dnsproviders API
type DNSProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSProviderSpec   `json:"spec,omitempty"`
	Status DNSProviderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DNSProviderList contains a list of DNSProvider
type DNSProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSProvider `json:"items"`
}

// GetProviderType returns a string representing the type of the provider represented by this resource.
func (resource *DNSProvider) GetProviderType() (string, error) {
	if resource.Spec.Dummy != nil {
		return "dummy", nil
	} else if resource.Spec.RFC2136 != nil {
		return "rfc2136", nil
	} else if resource.Spec.Cloudflare != nil {
		return "cloudflare", nil
	} else {
		return "", fmt.Errorf("Unknown provider type")
	}
}

func init() {
	SchemeBuilder.Register(&DNSProvider{}, &DNSProviderList{})
}
