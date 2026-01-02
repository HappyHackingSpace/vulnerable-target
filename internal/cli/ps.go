package cli

import (
	"os"
	"time"

	"github.com/happyhackingspace/vulnerable-target/internal/state"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/registry"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/disk"
	tmpl "github.com/happyhackingspace/vulnerable-target/pkg/template"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running deployments and their status",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := disk.NewConfig().WithFileName("deployments.db").WithBucketName("deployments")
		st, err := state.NewManager(cfg)
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		deployments, err := st.ListDeployments()
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleDefault)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Provider Name", "Template ID", "Status", "Created At"})

		count := 0
		for _, deployment := range deployments {
			provider := registry.GetProvider(deployment.ProviderName)
			if provider == nil {
				log.Error().Msgf("provider %q not found", deployment.ProviderName)
				continue
			}
			template, err := tmpl.GetByID(deployment.TemplateID)
			if err != nil {
				log.Error().Msgf("%v", err)
				continue
			}

			status := "unknown"
			if s, err := provider.Status(template); err != nil {
				log.Error().Msgf("%v", err)
			} else {
				status = s
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
