package providers

import (
	"fmt"
	"sync"

	dnsv1alpha1 "github.com/95ulisse/dns-operator/pkg/api/v1alpha1"
	types "github.com/95ulisse/dns-operator/pkg/types"
)

// ProviderContructor constructs a Provider given a kubernetes resource and a Context.
type ProviderContructor func(*types.ControllerContext, *dnsv1alpha1.DNSProvider) (types.Provider, error)

var (
	constructors     = make(map[string]ProviderContructor)
	constructorsLock sync.RWMutex
)

// RegisterProviderConstructor will register a provider constructor.
// `name` is a unique name representing the key of the corresponding k8s resource.
func RegisterProviderConstructor(name string, c ProviderContructor) {
	constructorsLock.Lock()
	defer constructorsLock.Unlock()
	constructors[name] = c
}

// ProviderFor builds a new Provider from the given kubernetes resource.
func ProviderFor(ctx *types.ControllerContext, resource *dnsv1alpha1.DNSProvider) (types.Provider, error) {
	providerType, err := resource.GetProviderType()
	if err != nil {
		return nil, fmt.Errorf("Could not get provider type: %s", err.Error())
	}

	constructorsLock.RLock()
	defer constructorsLock.RUnlock()
	if constructor, ok := constructors[providerType]; ok {
		return constructor(ctx, resource)
	}

	return nil, fmt.Errorf("Provider %s not registered", providerType)
}
