package cli

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	tmpl "github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop vulnerable environment by template id and provider",
	Run: func(cmd *cobra.Command, _ []string) {
		providerName, err := cmd.Flags().GetString("provider")
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		templateID, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		provider := registry.GetProvider(providerName)

		if provider == nil {
			log.Fatal().Msgf("provider %s not found", providerName)
		}

		template, err := tmpl.GetByID(templateID)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		err = provider.Stop(template)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		log.Info().Msgf("%s template stopped on %s", templateID, providerName)
	},
}

// setupStopCommand configures the stop command flags
func setupStopCommand() {
	stopCmd.Flags().StringP("provider", "p", "",
		fmt.Sprintf("Specify the provider for building a vulnerable environment (%s)",
			strings.Join(slices.Collect(maps.Keys(registry.Providers)), ", ")))

	stopCmd.Flags().String("id", "",
		"Specify a template ID for targeted vulnerable environment")

	err := stopCmd.MarkFlagRequired("provider")
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	err = stopCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
}
