package cli

import (
	"fmt"
	"strings"

	tmpl "github.com/happyhackingspace/vt/pkg/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newStopCommand creates the stop command.
func (c *CLI) newStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop vulnerable environment by template id or tags",
		Run: func(cmd *cobra.Command, _ []string) {
			providerName, err := cmd.Flags().GetString("provider")
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			templateID, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			tagsStr, err := cmd.Flags().GetString("tags")
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			if templateID == "" && tagsStr == "" {
				log.Fatal().Msg("either --id or --tags must be provided")
			}
			if templateID != "" && tagsStr != "" {
				log.Fatal().Msg("--id and --tags are mutually exclusive")
			}

			provider, ok := c.app.GetProvider(providerName)
			if !ok {
				log.Fatal().Msgf("provider %s not found", providerName)
			}

			if templateID != "" {
				template, err := tmpl.GetByID(c.app.Templates, templateID)
				if err != nil {
					log.Fatal().Msgf("%v", err)
				}

				err = provider.Stop(template)
				if err != nil {
					log.Fatal().Msgf("%v", err)
				}

				log.Info().Msgf("%s template stopped on %s", templateID, providerName)
			} else {
				tags := strings.Split(tagsStr, ",")
				templates, err := tmpl.GetByTags(c.app.Templates, tags)
				if err != nil {
					log.Fatal().Msgf("%v", err)
				}

				log.Info().Msgf("found %d templates matching tags: %s", len(templates), tagsStr)

				var failed []string
				var stopped int
				for _, template := range templates {
					if err := provider.Stop(template); err != nil {
						log.Error().Msgf("failed to stop %s: %v", template.ID, err)
						failed = append(failed, template.ID)
						continue
					}
					stopped++
					log.Info().Msgf("%s template stopped on %s", template.ID, providerName)
				}

				if len(failed) > 0 {
					log.Warn().Msgf("failed to stop %d templates: %s", len(failed), strings.Join(failed, ", "))
				}
				if stopped > 0 {
					log.Info().Msgf("successfully stopped %d templates", stopped)
				}
			}
		},
	}

	cmd.Flags().StringP("provider", "p", "docker-compose",
		fmt.Sprintf("Specify the provider for building a vulnerable environment (%s)",
			strings.Join(c.providerNames(), ", ")))

	cmd.Flags().String("id", "",
		"Specify a template ID for targeted vulnerable environment")

	cmd.Flags().StringP("tags", "t", "",
		"Specify comma-separated tags to stop all matching templates (e.g., --tags sqli,xss)")

	return cmd
}
