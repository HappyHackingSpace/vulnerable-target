package cli

import (
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
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
			templates.ListWithFilter(filter)
			return
		}

		if update {
			templates.Update()
			return
		}
	},
}

func setupTemplateCommand() {
	templateCmd.Flags().BoolP("list", "l", false, "List available templates")
	templateCmd.Flags().BoolP("update", "u", false, "Fetch templates repository to local working directory")
	templateCmd.Flags().StringP("filter", "f", "", "Filter templates by tag or keyword (only works with --list)")
}
