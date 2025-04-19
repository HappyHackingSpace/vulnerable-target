package main

import (
	"github.com/happyhackingspace/vulnerable-target/internal/cli"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/pkg/options"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/rs/zerolog/log"
)

func init() {
	logger.Init()
	templates.Init()
	options.LoadEnv()
}

func main() {
	cli.Execute()
	options := options.GetOptions()
	provider := registry.GetProvider(options.ProviderName)
	if provider == nil {
		log.Fatal().Msgf("provider %s not found", options.ProviderName)
	}
	if err := provider.Start(); err != nil {
		log.Fatal().Msgf("%v", err)
	}
	log.Info().Msgf("%s template is running on %s", options.TemplateID, options.ProviderName)
}
