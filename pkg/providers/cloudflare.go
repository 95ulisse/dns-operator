package providers

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	"github.com/95ulisse/dns-operator/pkg/types"
)

const (
	proxiedAnnotation string = "dns.k8s.marcocameriero.net/cloudflare-proxied"
)

// Cloudflare DNS provider.
type Cloudflare struct {
	log                logr.Logger
	zones              []dnsname.Name
	cf                 *cloudflare.API
	proxiedByDefault   bool
	cfZonesIDCache     map[string]string
	cfZonesIDCacheLock sync.RWMutex
}

// NewCloudflare creates a new instance of the Cloudflare provider.
func NewCloudflare(log logr.Logger, zones []dnsname.Name, cf *cloudflare.API, proxiedByDefault bool) *Cloudflare {
	return &Cloudflare{
		log:              log.WithName("providers").WithName("Cloudflare"),
		zones:            zones,
		cf:               cf,
		proxiedByDefault: proxiedByDefault,
		cfZonesIDCache:   make(map[string]string),
	}
}

// Zones returns a slice containing the DNS zones managed by this provider.
func (cf *Cloudflare) Zones() []dnsname.Name {
	return cf.zones
}

// UpdateRecord reconciles the given RRset with the records registered on Cloudflare.
func (cf *Cloudflare) UpdateRecord(zone dnsname.Name, resource v1alpha1.DNSRecord) error {

	// Retrieve the list of records of the RRset already registered on Cloudflare
	zoneID, err := cf.zoneIDFromName(zone)
	if err != nil {
		return err
	}
	filter := cloudflare.DNSRecord{
		Type: resource.RType(),
		Name: resource.Spec.Name.String(),
	}
	recordsAlreadyPresent, err := cf.cf.DNSRecords(zoneID, filter)
	if err != nil {
		return err
	}

	// Perform a diff between the wanted and the present records
	var toCreate []cloudflare.DNSRecord
	var toRemove []string
	toCreate, err = cf.toCFRecords(&resource)
	if err != nil {
		return err
	}
	for _, rrAlreadyPresent := range recordsAlreadyPresent {

		// If this record already present on CF is wanted by the user keep it, otherwise delete it
		found := -1
		for i, wanted := range toCreate {
			if rrEquals(&wanted, &rrAlreadyPresent) {
				found = i
				break
			}
		}

		if found >= 0 {
			toCreate = removeRR(toCreate, found)
		} else {
			toRemove = append(toRemove, rrAlreadyPresent.ID)
		}

	}

	// Synchronize the diff with cloudflare
	for _, id := range toRemove {
		cf.log.V(1).Info("Deleting old DNS record", "id", id)
		err := cf.cf.DeleteDNSRecord(zoneID, id)
		if err != nil {
			return err
		}
	}
	for _, rr := range toCreate {
		cf.log.V(1).Info("Creating new DNS record", "record", rr)
		_, err := cf.cf.CreateDNSRecord(zoneID, rr)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRecord deletes the given RRset from Cloudflare.
func (cf *Cloudflare) DeleteRecord(zone dnsname.Name, resource v1alpha1.DNSRecord) error {

	// Retrieve the list of records of the RRset already registered on Cloudflare
	zoneID, err := cf.zoneIDFromName(zone)
	if err != nil {
		return err
	}
	filter := cloudflare.DNSRecord{
		Type: resource.RType(),
		Name: resource.Spec.Name.String(),
	}
	recordsAlreadyPresent, err := cf.cf.DNSRecords(zoneID, filter)
	if err != nil {
		return err
	}

	// Delete all the records
	for _, rr := range recordsAlreadyPresent {
		cf.log.V(1).Info("Deleting old DNS record", "id", rr.ID)
		err := cf.cf.DeleteDNSRecord(zoneID, rr.ID)
		if err != nil {
			return err
		}
	}

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

// toCFRecords converts a DNSRecord resource (which represent a whole rrset) to a slice of Cloudflare records.
func (cf *Cloudflare) toCFRecords(resource *v1alpha1.DNSRecord) ([]cloudflare.DNSRecord, error) {

	// TTL
	var ttl = 1
	if resource.Spec.TTLSeconds != nil {
		ttl = int(*resource.Spec.TTLSeconds)
	}

	// Proxied attribute can be overridden with an annotation
	proxied := cf.proxiedByDefault
	if proxiedOverride, ok := resource.ObjectMeta.Annotations[proxiedAnnotation]; ok {
		if b, err := strconv.ParseBool(proxiedOverride); err == nil {
			proxied = b
		}
	}

	rrset := make([]cloudflare.DNSRecord, 0, 1)
	push := func(content string, priority int) {
		var rr cloudflare.DNSRecord
		rr.Type = resource.RType()
		rr.Name = resource.Spec.Name.String()
		rr.Content = content
		rr.TTL = ttl
		rr.Priority = priority
		rr.Proxied = proxied
		rrset = append(rrset, rr)
	}

	switch resource.RType() {
	case "A":
		for _, value := range resource.Spec.RRSet.A {
			push(value.String(), 0)
		}

	case "AAAA":
		for _, value := range resource.Spec.RRSet.AAAA {
			push(value.String(), 0)
		}

	case "CNAME":
		for _, value := range resource.Spec.RRSet.CNAME {
			push(value.String(), 0)
		}

	case "TXT":
		for _, value := range resource.Spec.RRSet.TXT {
			push(value, 0)
		}

	case "MX":
		for _, mx := range resource.Spec.RRSet.MX {
			push(mx.Host.String(), int(mx.Preference))
		}

	default:
		return nil, fmt.Errorf("Unsupported DNS record")
	}

	return rrset, nil
}

func rrEquals(rr1 *cloudflare.DNSRecord, rr2 *cloudflare.DNSRecord) bool {
	return rr1.Type == rr2.Type &&
		rr1.Name == rr2.Name &&
		rr1.Content == rr2.Content &&
		rr1.Proxied == rr2.Proxied &&
		rr1.TTL == rr2.TTL &&
		rr1.Priority == rr2.Priority
}

// removeRR will remove the item with index `i` from the given slice.
// NOTE: this function does not preserve the order of the original slice.
func removeRR(s []cloudflare.DNSRecord, i int) []cloudflare.DNSRecord {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
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

		var proxiedByDefault = true
		if resource.Spec.Cloudflare.ProxiedByDefault != nil {
			proxiedByDefault = *resource.Spec.Cloudflare.ProxiedByDefault
		}

		return NewCloudflare(ctx.Log, resource.Spec.Zones, cf, proxiedByDefault), nil
	})
}
