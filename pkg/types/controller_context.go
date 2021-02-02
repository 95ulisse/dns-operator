package types

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ControllerContext contains structures shared with all the controllers of the application.
type ControllerContext struct {
	RootContext   context.Context
	Client        client.Client
	Log           logr.Logger
	EventRecorder record.EventRecorder

	providers     map[string]Provider
	providersLock sync.RWMutex
}

// GetProvider returns the provider with the registered name, if any.
func (ctx *ControllerContext) GetProvider(name string, provider *Provider) bool {
	ctx.providersLock.RLock()
	defer ctx.providersLock.RUnlock()

	if ctx.providers == nil {
		return false
	}

	p, present := ctx.providers[name]
	if present {
		*provider = p
	}
	return present
}

// RemoveProvider removed the provider with the given name from the global map of registered providers.
func (ctx *ControllerContext) RemoveProvider(name string) {
	ctx.providersLock.Lock()
	defer ctx.providersLock.Unlock()

	if ctx.providers != nil {
		delete(ctx.providers, name)
	}
}

// SetProvider registers a new provider with the given name.
func (ctx *ControllerContext) SetProvider(name string, provider Provider) {
	ctx.providersLock.Lock()
	defer ctx.providersLock.Unlock()

	if ctx.providers == nil {
		ctx.providers = make(map[string]Provider)
	}

	ctx.providers[name] = provider
}
