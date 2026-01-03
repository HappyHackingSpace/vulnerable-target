package cli

import (
	"fmt"

	templ "github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (c *CLI) newInspectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "inspect operations",
		Run: func(cmd *cobra.Command, _ []string) {
			templateID, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			template, err := templ.GetByID(c.app.Templates, templateID)
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			fmt.Println(template.String())
		},
	}
	cmd.Flags().String("id", "", "Specify a template ID for targeted vulnerable environment")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		panic(err)
	}
	return cmd
}
