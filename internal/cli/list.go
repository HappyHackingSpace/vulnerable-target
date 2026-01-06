package cli

import (
	tmpl "github.com/happyhackingspace/vt/pkg/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// newTemplateCommand creates the template command.
func (c *CLI) newTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "template operations",
		Run: func(cmd *cobra.Command, _ []string) {
			list, err := cmd.Flags().GetBool("list")
			if err != nil {
				log.Error().Err(err).Msg("failed to get list flag")
				return
			}

			update, err := cmd.Flags().GetBool("update")
			if err != nil {
				log.Error().Err(err).Msg("failed to get update flag")
				return
			}

			filter, err := cmd.Flags().GetString("filter")
			if err != nil {
				log.Error().Err(err).Msg("failed to get filter flag")
				return
			}

			if list && update {
				log.Error().Msg("only one of --list, or --update can be specified")
				return
			}

			if filter != "" && !list {
				log.Error().Msg("--filter can only be used with --list")
				return
			}

			if !list && !update {
				if err := cmd.Help(); err != nil {
					log.Error().Err(err).Msg("failed to display help")
				}
				return
			}

			if list {
				tmpl.ListTemplatesWithFilter(c.app.Templates, filter)
				return
			}

			if update {
				if err := tmpl.SyncTemplates(c.app.Config.TemplatesPath); err != nil {
					log.Error().Err(err).Msg("failed to sync templates")
					return
				}
				// Reload templates after sync
				templates, err := tmpl.LoadTemplates(c.app.Config.TemplatesPath)
				if err != nil {
					log.Error().Err(err).Msg("failed to reload templates")
					return
				}
				// Update the app's templates
				c.app.Templates = templates
				log.Info().Msg("Templates updated successfully")
				return
			}
		},
	}

	cmd.Flags().BoolP("list", "l", false, "List available templates")
	cmd.Flags().BoolP("update", "u", false, "Fetch templates repository to local working directory")
	cmd.Flags().StringP("filter", "f", "", "Filter templates by tag or keyword (only works with --list)")

	return cmd
}
