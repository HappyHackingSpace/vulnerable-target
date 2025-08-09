package cli

import (
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/spf13/cobra"
)

var (
	filterTag string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates with descriptions",
	Run: func(_ *cobra.Command, _ []string) {
		templates.ListWithFilter(filterTag)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&filterTag, "filter", "f", "", "Filter templates by tag (e.g., --filter=php or -f php)")
}
