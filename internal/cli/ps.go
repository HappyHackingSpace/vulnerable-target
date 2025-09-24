package cli

import (
	"os"
	"time"

	"github.com/happyhackingspace/vulnerable-target/internal/state"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Runs selected template on chosen provider",
	Run: func(_ *cobra.Command, _ []string) {
		st, err := state.NewManager()
		if err != nil {
			log.Error().Msgf("%v", err)

		}

		deployments, err := st.ListDeployments()
		if err != nil {
			log.Error().Msgf("%v", err)
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleDefault)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Provider Name", "Template ID", "Status", "Created At"})

		count := 0
		for _, deployment := range deployments {
			provider := registry.GetProvider(deployment.ProviderName)
			template, err := templates.GetByID(deployment.TemplateID)
			if err != nil {
				log.Error().Msgf("%v", err)
			}

			status, err := provider.Status(template)
			if err != nil {
				log.Error().Msgf("%v", err)
			}

			t.AppendRow(table.Row{
				deployment.ProviderName,
				deployment.TemplateID,
				status,
				deployment.CreatedAt.Format(time.DateTime),
			})
			count++
		}

		if count == 0 {
			log.Info().Msg("there is no running environment")
			return
		}

		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
