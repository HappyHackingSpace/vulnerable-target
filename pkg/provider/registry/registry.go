// Package registry manages the collection of available providers.
package registry

import (
	"github.com/happyhackingspace/vt/internal/state"
	"github.com/happyhackingspace/vt/pkg/provider"
	"github.com/happyhackingspace/vt/pkg/provider/dockercompose"
)

// NewProviders creates and returns a map of all available providers.
// Each provider is initialized with the given state manager.
func NewProviders(sm *state.Manager) map[string]provider.Provider {
	return map[string]provider.Provider{
		"docker-compose": dockercompose.NewDockerCompose(sm),
	}
}

// GetProvider retrieves a provider by name from the given providers map.
func GetProvider(providers map[string]provider.Provider, name string) (provider.Provider, bool) {
	p, ok := providers[name]
	return p, ok
}
