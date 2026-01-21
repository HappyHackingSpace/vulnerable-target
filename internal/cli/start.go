package cli

import (
	"fmt"
	"strings"

	tmpl "github.com/happyhackingspace/vt/pkg/template"
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

				c.startTemplate(provider, template, providerName)
			} else {
				tags := strings.Split(tagsStr, ",")
				templates, err := tmpl.GetByTags(c.app.Templates, tags)
				if err != nil {
					log.Fatal().Msgf("%v", err)
				}

				log.Info().Msgf("found %d templates matching tags: %s", len(templates), tagsStr)

				var failed []string
				for _, template := range templates {
					if err := provider.Start(template); err != nil {
						log.Error().Msgf("failed to start %s: %v", template.ID, err)
						failed = append(failed, template.ID)
						continue
					}

					if len(template.PostInstall) > 0 {
						log.Info().Msgf("Post-installation instructions for %s:", template.ID)
						for _, instruction := range template.PostInstall {
							fmt.Printf("  %s\n", instruction)
						}
					}

					log.Info().Msgf("%s template is running on %s", template.ID, providerName)
				}

				if len(failed) > 0 {
					log.Warn().Msgf("failed to start %d templates: %s", len(failed), strings.Join(failed, ", "))
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
		"Specify comma-separated tags to start all matching templates (e.g., --tags sqli,xss)")

	return cmd
}

// startTemplate starts a single template and logs the result.
func (c *CLI) startTemplate(provider interface{ Start(*tmpl.Template) error }, template *tmpl.Template, providerName string) {
	err := provider.Start(template)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	if len(template.PostInstall) > 0 {
		log.Info().Msg("Post-installation instructions:")
		for _, instruction := range template.PostInstall {
			fmt.Printf("  %s\n", instruction)
		}
	}

	log.Info().Msgf("%s template is running on %s", template.ID, providerName)
}
