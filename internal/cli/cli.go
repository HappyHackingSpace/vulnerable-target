// Package cli provides command-line interface functionality for the vulnerable target application.
package cli

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/happyhackingspace/vulnerable-target/internal/app"
	"github.com/happyhackingspace/vulnerable-target/internal/banner"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// logLevels defines the valid log levels supported by the application.
var logLevels = map[string]bool{
	zerolog.DebugLevel.String(): true,
	zerolog.InfoLevel.String():  true,
	zerolog.WarnLevel.String():  true,
	zerolog.ErrorLevel.String(): true,
	zerolog.FatalLevel.String(): true,
	zerolog.PanicLevel.String(): true,
}

// CLI encapsulates the command-line interface with its dependencies.
type CLI struct {
	app     *app.App
	rootCmd *cobra.Command
}

// New creates a new CLI instance with the given application context.
func New(application *app.App) *CLI {
	c := &CLI{
		app: application,
	}
	c.setupCommands()
	return c
}

// setupCommands initializes all CLI commands and their configurations.
func (c *CLI) setupCommands() {
	c.rootCmd = &cobra.Command{
		Use:     "vt",
		Short:   "Create vulnerable environment",
		Version: banner.AppVersion,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			verbosityLevel, err := cmd.Flags().GetString("verbosity")
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
			logger.InitWithLevel(verbosityLevel)
		},
		SilenceErrors: true,
	}
	c.rootCmd.SetHelpTemplate(banner.Banner() + "\n" + c.rootCmd.HelpTemplate())

	// Setup root flags
	c.rootCmd.PersistentFlags().StringP("verbosity", "v", zerolog.InfoLevel.String(),
		fmt.Sprintf("Set the verbosity level for logs (%s)",
			strings.Join(slices.Collect(maps.Keys(logLevels)), ", ")))

	// Register all subcommands
	c.rootCmd.AddCommand(c.newStartCommand())
	c.rootCmd.AddCommand(c.newStopCommand())
	c.rootCmd.AddCommand(c.newPsCommand())
	c.rootCmd.AddCommand(c.newTemplateCommand())
	c.rootCmd.AddCommand(c.newInspectCommand())
}

// Run executes the CLI and returns any error.
func (c *CLI) Run() error {
	return c.rootCmd.Execute()
}

// providerNames returns a slice of registered provider names.
func (c *CLI) providerNames() []string {
	return slices.Collect(maps.Keys(c.app.Providers))
}
