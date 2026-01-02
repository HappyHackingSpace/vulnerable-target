package cli

import (
	"fmt"
	"strings"

	tmpl "github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newStartCommand creates the start command.
func (c *CLI) newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
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

			if len(templateID) == 0 {
				if err := cmd.Help(); err != nil {
					log.Fatal().Msgf("%v", err)
				}
				return
			}

			provider, ok := c.app.GetProvider(providerName)
			if !ok {
				log.Fatal().Msgf("provider %s not found", providerName)
			}

			template, err := tmpl.GetByID(c.app.Templates, templateID)
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

	cmd.Flags().StringP("provider", "p", "docker-compose",
		fmt.Sprintf("Specify the provider for building a vulnerable environment (%s)",
			strings.Join(c.providerNames(), ", ")))

	cmd.Flags().String("id", "",
		"Specify a template ID for targeted vulnerable environment")

	if err := cmd.MarkFlagRequired("id"); err != nil {
		log.Fatal().Msgf("%v", err)
	}

	return cmd
}
