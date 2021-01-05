package providers

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/miekg/dns"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/api/v1alpha1"
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
	zone       string
	useTsig    bool
	keyName    string
	algorithm  string
	nameserver string
}

// UpdateRecord updates a record on the backend server.
func (provider *RFC2136) UpdateRecord(record *dnsv1alpha1.DNSRecord) error {

	// Prepare the DNS message
	rr, err := extractRR(record)
	if err != nil {
		return err
	}
	msg := new(dns.Msg)
	msg.SetUpdate(provider.zone)
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

	provider.log.Info("Update successful")

	return nil
}

// DeleteRecord deletes a record from the backend server.
func (provider *RFC2136) DeleteRecord(record *dnsv1alpha1.DNSRecord) error {

	// Prepare the DNS message
	rr, err := extractRR(record)
	if err != nil {
		return err
	}
	msg := new(dns.Msg)
	msg.SetUpdate(provider.zone)
	msg.Remove([]dns.RR{rr})
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

	provider.log.Info("Delete successful")

	return nil
}

func extractTSIGKey(resource *dnsv1alpha1.DNSProvider, dnsClient *dns.Client, k8sClient client.Client) (string, string, error) {

	// Extract the required parameters
	secretRef := resource.Spec.RFC2136.TSIGSecretRef
	keyName := resource.Spec.RFC2136.TSIGKeyName
	algorithm := resource.Spec.RFC2136.TSIGAlgorithm
	if secretRef == nil || keyName == nil || algorithm == nil {
		return "", "", fmt.Errorf("All fields tsigSecretRef, tsigKeyName and tsigAlgorithm are required when specifying a TSIG key")
	}
	if !strings.HasSuffix(*keyName, ".") {
		appended := fmt.Sprintf("%s%s", *keyName, ".")
		keyName = &appended
	}

	// Check that the algorithm name is valid
	dnsAlgorithm, algorithmSupported := supportedAlgorithms[*algorithm]
	if !algorithmSupported {
		return "", "", fmt.Errorf("Unsupported TSIG key algorithm %s", *algorithm)
	}

	// Resolve the secret reference
	secretName := secretRef.Name
	secretNamespace := secretRef.Namespace
	if secretNamespace == nil {
		secretNamespace = &resource.Namespace
	}
	var secret corev1.Secret
	if err := k8sClient.Get(context.Background(), types.NamespacedName{Name: secretName, Namespace: *secretNamespace}, &secret); err != nil {
		return "", "", err
	}

	// Extract the key from the secret
	key, keyPresent := secret.Data[secretRef.Key]
	if !keyPresent {
		return "", "", fmt.Errorf("Cannot find key %s in secret %s/%s", secretRef.Key, *secretNamespace, secretName)
	}

	// Configure the dns client
	dnsClient.TsigSecret = make(map[string]string)
	dnsClient.TsigSecret[*keyName] = string(key)

	return *keyName, dnsAlgorithm, nil

}

func extractRR(record *dnsv1alpha1.DNSRecord) (dns.RR, error) {

	var name string = dns.Fqdn(record.Spec.Name)
	var ttl uint32 = 3600
	if record.Spec.TTLSeconds != nil {
		ttl = *record.Spec.TTLSeconds
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
		rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
		rr.A = ip
		return rr, nil
	}

	return nil, fmt.Errorf("Unsupported DNS record")

}

func init() {
	RegisterProviderFactory(func(resource *dnsv1alpha1.DNSProvider, k8sClient client.Client, log logr.Logger) (Provider, error) {

		// Skip this provider if not configured
		if resource.Spec.RFC2136 == nil {
			return nil, nil
		}

		// Prepare a DNS client pre-configured with the TSIG secrets
		var keyName, algorithm string
		var useTsig = resource.Spec.RFC2136.TSIGSecretRef != nil
		dnsClient := new(dns.Client)
		if useTsig {
			var err error
			if keyName, algorithm, err = extractTSIGKey(resource, dnsClient, k8sClient); err != nil {
				return nil, err
			}
		}

		rfc2136 := RFC2136{
			client:     dnsClient,
			zone:       dns.Fqdn(resource.Spec.Zone),
			useTsig:    useTsig,
			keyName:    keyName,
			algorithm:  algorithm,
			nameserver: resource.Spec.RFC2136.Nameserver,
			log: log.
				WithName("providers").WithName("RFC2136").
				WithValues("dnsprovider", fmt.Sprintf("%s/%s", resource.Namespace, resource.Name)),
		}
		return &rfc2136, nil

	})
}
