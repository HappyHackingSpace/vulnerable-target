// Package app provides the application dependency container and configuration.
package app

import (
	"os"
	"path/filepath"

	"github.com/happyhackingspace/vulnerable-target/internal/state"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider"
	"github.com/happyhackingspace/vulnerable-target/pkg/template"
)

// Config holds application configuration.
type Config struct {
	TemplatesPath string
	StoragePath   string
	LogLevel      string
}

// App is the dependency container for the application.
type App struct {
	Templates    map[string]template.Template
	Providers    map[string]provider.Provider
	StateManager *state.Manager
	Config       *Config
}

// DefaultConfig returns the default application configuration.
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Config{
		TemplatesPath: filepath.Join(homeDir, "vt-templates"),
		StoragePath:   filepath.Join(homeDir, ".vt-cli"),
		LogLevel:      "info",
	}
}

// NewApp creates a new App instance with the given dependencies.
func NewApp(
	templates map[string]template.Template,
	providers map[string]provider.Provider,
	stateManager *state.Manager,
	config *Config,
) *App {
	return &App{
		Templates:    templates,
		Providers:    providers,
		StateManager: stateManager,
		Config:       config,
	}
}

// GetProvider retrieves a provider by name from the app's provider map.
func (a *App) GetProvider(name string) (provider.Provider, bool) {
	p, ok := a.Providers[name]
	return p, ok
}
