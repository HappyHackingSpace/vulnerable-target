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

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Runs selected template on chosen provider",
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
		if len(templateID) == 0 {
			err := cmd.Help()
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
			return
		}

		if provider == nil {
			log.Fatal().Msgf("provider %s not found", providerName)
		}

		template, err := tmpl.GetByID(templateID)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		err = provider.Start(template)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		if len(template.PostInstall) > 0 {

			log.Info().Msg("Post-installation instructions:")
			for _, instruction := range template.PostInstall {
				fmt.Printf("  %s\n", instruction)
			}
		}

		log.Info().Msgf("%s template is running on %s", templateID, providerName)
	},
}

// setupStartCommand configures the start command flags
func setupStartCommand() {
	startCmd.Flags().StringP("provider", "p", "docker-compose",
		fmt.Sprintf("Specify the provider for building a vulnerable environment (%s)",
			strings.Join(slices.Collect(maps.Keys(registry.Providers)), ", ")))

	startCmd.Flags().String("id", "",
		"Specify a template ID for targeted vulnerable environment")

	err := startCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
}
