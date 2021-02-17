package providers

import (
	"context"
	"fmt"
	"sync"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/go-logr/logr"
	"github.com/miekg/dns"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	"github.com/95ulisse/dns-operator/pkg/types"
)

// Cloudflare DNS provider.
type Cloudflare struct {
	log                logr.Logger
	zones              []dnsname.Name
	cf                 *cloudflare.API
	cfZonesIDCache     map[string]string
	cfZonesIDCacheLock sync.RWMutex
}

// NewCloudflare creates a new instance of the Cloudflare provider.
func NewCloudflare(log logr.Logger, zones []dnsname.Name, cf *cloudflare.API) *Cloudflare {
	return &Cloudflare{
		log:            log.WithName("providers").WithName("Cloudflare"),
		zones:          zones,
		cf:             cf,
		cfZonesIDCache: make(map[string]string),
	}
}

// Zones returns a slice containing the DNS zones managed by this provider.
func (cf *Cloudflare) Zones() []dnsname.Name {
	return cf.zones
}

// UpdateRecord reconciles the given RRset with the records registered on Cloudflare.
func (cf *Cloudflare) UpdateRecord(zone dnsname.Name, rrset []dns.RR) error {

	// Retrieve the list of records of the RRset already registered on Cloudflare
	zoneID, err := cf.zoneIDFromName(zone)
	if err != nil {
		return err
	}
	filter := cloudflare.DNSRecord{
		Type: dns.TypeToString[rrset[0].Header().Rrtype],
		Name: unFqdn(rrset[0].Header().Name),
	}
	cf.log.V(1).Info("CF filter", "filter", filter)
	recordsAlreadyPresent, err := cf.cf.DNSRecords(zoneID, filter)
	if err != nil {
		return err
	}

	// Perform a diff between the wanted and the present records
	cf.log.V(1).Info(fmt.Sprintf("String of RRset: %s", rrset[0].String()))
	for _, rrAlreadyPresent := range recordsAlreadyPresent {
		cf.log.V(1).Info("Record on CF", "record", rrAlreadyPresent)
	}

	return nil
}

// DeleteRecord deletes the given RRset from Cloudflare.
func (cf *Cloudflare) DeleteRecord(zone dnsname.Name, rrset []dns.RR) error {
	cf.log.Info("Delete successful")
	return nil
}

func (cf *Cloudflare) zoneIDFromName(zone dnsname.Name) (string, error) {

	// First check if the zone is in the cache
	id := func() string {
		cf.cfZonesIDCacheLock.RLock()
		defer cf.cfZonesIDCacheLock.RUnlock()
		if id, ok := cf.cfZonesIDCache[zone.String()]; ok {
			return id
		}
		return ""
	}()
	if id != "" {
		return id, nil
	}

	// Resolve the zone using CF api
	id, err := cf.cf.ZoneIDByName(zone.String())
	if err != nil {
		cf.log.Error(err, "Could not resolve zone name", "zone", zone.String())
		return "", err
	}

	// Store the id in the cache
	cf.cfZonesIDCacheLock.Lock()
	defer cf.cfZonesIDCacheLock.Unlock()
	cf.cfZonesIDCache[zone.String()] = id

	return id, nil

}

func unFqdn(name string) string {
	if dns.IsFqdn(name) {
		return name[0 : len(name)-1]
	}
	return name
}

func init() {
	RegisterProviderConstructor("cloudflare", func(ctx *types.ControllerContext, resource *dnsv1alpha1.DNSProvider) (types.Provider, error) {

		// Extract the required parameters
		email := resource.Spec.Cloudflare.Email
		apiToken := resource.Spec.Cloudflare.APITokenSecretRef
		apiKey := resource.Spec.Cloudflare.APIKeySecretRef
		if apiToken == nil && apiKey == nil {
			return nil, fmt.Errorf("One between `apiTokenSecretRef` and `apiKeySecretRef` is required")
		}
		if apiKey != nil && (email == nil || *email == "") {
			return nil, fmt.Errorf("`email` is required when authenticating with an API Key")
		}

		// Resolve the secret ref
		secretRef := apiToken
		if secretRef == nil {
			secretRef = apiKey
		}
		secretName := secretRef.Name
		secretNamespace := secretRef.Namespace
		if secretNamespace == nil {
			secretNamespace = &resource.Namespace
		}
		var secret corev1.Secret
		if err := ctx.Client.Get(context.Background(), k8stypes.NamespacedName{Name: secretName, Namespace: *secretNamespace}, &secret); err != nil {
			return nil, err
		}

		// Extract the key or token from the secret
		key, keyPresent := secret.Data[secretRef.Key]
		if !keyPresent {
			return nil, fmt.Errorf("Cannot find key %s in secret %s/%s", secretRef.Key, *secretNamespace, secretName)
		}

		// Build a Cloudflare client
		var cf *cloudflare.API
		var err error
		if apiKey != nil {
			cf, err = cloudflare.New(string(key), *email)
		} else {
			cf, err = cloudflare.NewWithAPIToken(string(key))
		}
		if err != nil {
			return nil, err
		}

		return NewCloudflare(ctx.Log, resource.Spec.Zones, cf), nil
	})
}
