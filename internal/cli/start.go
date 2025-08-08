package cli

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/pkg/options"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, _ []string) {
		options := options.GetOptions()
		provider := registry.GetProvider(options.ProviderName)
		if len(options.TemplateID) == 0 || len(options.ProviderName) == 0 {
			err := cmd.Help()
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
			return
		}
		if provider == nil {
			log.Fatal().Msgf("provider %s not found", options.ProviderName)
		}

		template, err := templates.GetByID(options.TemplateID)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		err = provider.Start(template)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		log.Info().Msgf("%s template is running on %s", options.TemplateID, options.ProviderName)
	},
}

func init() {
	options := options.GetOptions()
	logger.Init()
	templates.Init()

	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&options.ProviderName, "provider", "p", "",
		fmt.Sprintf("Specify the provider for building a vulnerable environment (%s)",
			strings.Join(slices.Collect(maps.Keys(registry.Providers)), ", ")))

	startCmd.Flags().StringVar(&options.TemplateID, "id", "",
		"Specify a template ID for targeted vulnerable environment")
}
