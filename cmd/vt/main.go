// Package main is the entry point for the vulnerable target application.
package main

import (
	"os"

	"github.com/happyhackingspace/vulnerable-target/internal/app"
	"github.com/happyhackingspace/vulnerable-target/internal/banner"
	"github.com/happyhackingspace/vulnerable-target/internal/cli"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/internal/state"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/disk"
	"github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize configuration
	cfg := app.DefaultConfig()

	// Initialize logger (sets global logger for zerolog)
	appLogger := logger.NewWithLevel(cfg.LogLevel)
	logger.SetGlobal(appLogger)

	// Print banner
	banner.Print()

	// Load templates from repository
	templates, err := template.LoadTemplates(cfg.TemplatesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load templates")
	}

	// Create state manager
	storeCfg := disk.NewConfig().
		WithFileName("deployments.db").
		WithBucketName("deployments")
	stateManager, err := state.NewManager(storeCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create state manager")
	}

	// Create providers with injected dependencies
	providers := registry.NewProviders(stateManager)

	// Create application context
	application := app.NewApp(templates, providers, stateManager, cfg)

	// Create and run CLI
	if err := cli.New(application).Run(); err != nil {
		log.Fatal().Err(err).Msg("CLI error")
		os.Exit(1)
	}
}
