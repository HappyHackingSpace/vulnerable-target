// Package main is the entry point for the vulnerable target application.
package main

import (
	"github.com/happyhackingspace/vulnerable-target/internal/app"
	"github.com/happyhackingspace/vulnerable-target/internal/cli"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/internal/state"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/disk"
	"github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := app.DefaultConfig()

	appLogger := logger.NewWithLevel(cfg.LogLevel)
	logger.SetGlobal(appLogger)

	templates, err := template.LoadTemplates(cfg.TemplatesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load templates")
	}

	storeCfg := disk.NewConfig().
		WithFileName("deployments.db").
		WithBucketName("deployments")
	stateManager, err := state.NewManager(storeCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create state manager")
	}

	providers := registry.NewProviders(stateManager)

	application := app.NewApp(templates, providers, stateManager, cfg)

	if err := cli.New(application).Run(); err != nil {
		log.Fatal().Err(err).Msg("CLI error")
	}
}
