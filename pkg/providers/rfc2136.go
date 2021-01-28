package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/miekg/dns"

	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	"github.com/95ulisse/dns-operator/pkg/types"
)

var supportedAlgorithms = map[string]string{
	"HMACMD5":    dns.HmacMD5,
	"HMACSHA1":   dns.HmacSHA1,
	"HMACSHA256": dns.HmacSHA256,
	"HMACSHA512": dns.HmacSHA512,
}

// RFC2136 is a DNS provider which uses Dynamic DNS (https://tools.ietf.org/html/rfc2136)
// for updates to a backend server.
type RFC2136 struct {
	log        logr.Logger
	client     *dns.Client
	zones      []dnsname.Name
	nameserver string
	useTsig    bool
	keyName    string
	algorithm  string
}

// NewRFC2136 creates a new DNS provider which uses Dynamic DNS for updates.
func NewRFC2136(log logr.Logger, zones []dnsname.Name, nameserver string) *RFC2136 {
	client := new(dns.Client)
	client.SingleInflight = true

	return &RFC2136{
		log:        log.WithName("providers").WithName("RFC2136"),
		client:     client,
		zones:      zones,
		nameserver: nameserver,
	}
}

// WithTsig configures transaction signatures for DNS updates.
func (provider *RFC2136) WithTsig(secret, keyName, algorithm string) *RFC2136 {
	provider.client.TsigSecret = make(map[string]string)
	provider.client.TsigSecret[keyName] = secret
	provider.keyName = keyName
	provider.algorithm = algorithm
	provider.useTsig = true
	return provider
}

// Zones returns a slice containing the DNS zones managed by this provider.
func (provider *RFC2136) Zones() []dnsname.Name {
	return provider.zones
}

// UpdateRecord updates a record on the backend server.
func (provider *RFC2136) UpdateRecord(zone dnsname.Name, rr dns.RR) error {

	// Prepare the DNS message
	msg := new(dns.Msg)
	msg.SetUpdate(zone.ToFQDN().String())
	msg.RemoveRRset([]dns.RR{rr})
	msg.Insert([]dns.RR{rr})
	if provider.useTsig {
		msg.SetTsig(provider.keyName, provider.algorithm, 300, time.Now().Unix())
	}

	// Send the message
	res, _, err := provider.client.Exchange(msg, provider.nameserver)
	if err != nil {
		return fmt.Errorf("DNS update failed: %s", err)
	}
	if res != nil && res.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("DNS update failed. Server replied: %s", dns.RcodeToString[res.Rcode])
	}

	provider.log.Info(fmt.Sprintf("Updated %s %s", dns.TypeToString[rr.Header().Rrtype], rr.Header().Name))

	return nil
}

// DeleteRecord deletes a record from the backend server.
func (provider *RFC2136) DeleteRecord(zone dnsname.Name, rr dns.RR) error {

	// Prepare the DNS message
	msg := new(dns.Msg)
	msg.SetUpdate(zone.ToFQDN().String())
	msg.RemoveRRset([]dns.RR{rr})
	if provider.useTsig {
		msg.SetTsig(provider.keyName, provider.algorithm, 300, time.Now().Unix())
	}

	// Send the message
	res, _, err := provider.client.Exchange(msg, provider.nameserver)
	if err != nil {
		return fmt.Errorf("DNS delete failed: %s", err)
	}
	if res != nil && res.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("DNS delete failed. Server replied: %s", dns.RcodeToString[res.Rcode])
	}

	provider.log.Info(fmt.Sprintf("Deleted %s %s", dns.TypeToString[rr.Header().Rrtype], rr.Header().Name))

	return nil
}

func extractTSIGKey(resource *dnsv1alpha1.DNSProvider, k8sClient client.Client) (string, string, string, error) {

	// Extract the required parameters
	secretRef := resource.Spec.RFC2136.TSIGSecretRef
	keyName := resource.Spec.RFC2136.TSIGKeyName
	algorithm := resource.Spec.RFC2136.TSIGAlgorithm
	if secretRef == nil || keyName == nil || algorithm == nil {
		return "", "", "", fmt.Errorf("All fields tsigSecretRef, tsigKeyName and tsigAlgorithm are required when specifying a TSIG key")
	}
	if !strings.HasSuffix(*keyName, ".") {
		appended := fmt.Sprintf("%s%s", *keyName, ".")
		keyName = &appended
	}

	// Check that the algorithm name is valid
	dnsAlgorithm, algorithmSupported := supportedAlgorithms[*algorithm]
	if !algorithmSupported {
		return "", "", "", fmt.Errorf("Unsupported TSIG key algorithm %s", *algorithm)
	}

	// Resolve the secret reference
	secretName := secretRef.Name
	secretNamespace := secretRef.Namespace
	if secretNamespace == nil {
		secretNamespace = &resource.Namespace
	}
	var secret corev1.Secret
	if err := k8sClient.Get(context.Background(), k8stypes.NamespacedName{Name: secretName, Namespace: *secretNamespace}, &secret); err != nil {
		return "", "", "", err
	}

	// Extract the key from the secret
	key, keyPresent := secret.Data[secretRef.Key]
	if !keyPresent {
		return "", "", "", fmt.Errorf("Cannot find key %s in secret %s/%s", secretRef.Key, *secretNamespace, secretName)
	}

	return *keyName, string(key), dnsAlgorithm, nil

}

func init() {
	RegisterProviderConstructor("rfc2136", func(ctx *types.ControllerContext, resource *dnsv1alpha1.DNSProvider) (types.Provider, error) {
		provider := NewRFC2136(ctx.Log, resource.Spec.Zones, resource.Spec.RFC2136.Nameserver)
		if resource.Spec.RFC2136.TSIGSecretRef != nil {
			keyName, secret, algorithm, err := extractTSIGKey(resource, ctx.Client)
			if err != nil {
				return nil, err
			}
			provider = provider.WithTsig(secret, keyName, algorithm)
		}

		return provider, nil
	})
}
