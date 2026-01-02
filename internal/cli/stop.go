package cli

import (
	"fmt"
	"strings"

	tmpl "github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newStopCommand creates the stop command.
func (c *CLI) newStopCommand() *cobra.Command {
	cmd := &cobra.Command{
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

			provider, ok := c.app.GetProvider(providerName)
			if !ok {
				log.Fatal().Msgf("provider %s not found", providerName)
			}

			template, err := tmpl.GetByID(c.app.Templates, templateID)
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

	cmd.Flags().StringP("provider", "p", "docker-compose",
		fmt.Sprintf("Specify the provider for building a vulnerable environment (%s)",
			strings.Join(c.providerNames(), ", ")))

	cmd.Flags().String("id", "",
		"Specify a template ID for targeted vulnerable environment")

	if err := cmd.MarkFlagRequired("provider"); err != nil {
		log.Fatal().Msgf("%v", err)
	}

	if err := cmd.MarkFlagRequired("id"); err != nil {
		log.Fatal().Msgf("%v", err)
	}

	return cmd
}
