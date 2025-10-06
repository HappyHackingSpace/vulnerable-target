package cli

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	banner "github.com/happyhackingspace/vulnerable-target/internal/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// LogLevels defines the valid log levels supported by the application.
var LogLevels = map[string]bool{
	zerolog.DebugLevel.String(): true,
	zerolog.InfoLevel.String():  true,
	zerolog.WarnLevel.String():  true,
	zerolog.ErrorLevel.String(): true,
	zerolog.FatalLevel.String(): true,
	zerolog.PanicLevel.String(): true,
}

// setupRootFlags configures the root command flags
func setupRootFlags() {
	rootCmd.PersistentFlags().StringP("verbosity", "v", zerolog.InfoLevel.String(),
		fmt.Sprintf("Set the verbosity level for logs (%s)",
			strings.Join(slices.Collect(maps.Keys(LogLevels)), ", ")))
}

var rootCmd = &cobra.Command{
	Use:     "vt",
	Short:   "Create vulnerable environment",
	Version: banner.AppVersion,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		verbosityLevel, err := cmd.Flags().GetString("verbosity")
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		logger.InitWithLevel(verbosityLevel)
		if cmd.Name() != "help" {
			banner.Print()
		}
	},
	SilenceErrors: true,
}

// InitCLI initializes all CLI commands and their configurations
func InitCLI() {
	// Setup root command flags
	setupRootFlags()

	// Register all subcommands
	registerCommands()
}

// registerCommands registers all CLI subcommands and configures their flags
func registerCommands() {
	// Register commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(validateCmd)

	// Setup command-specific flags
	setupListCommand()
	setupStartCommand()
	setupStopCommand()
	setupValidateCommand()
}

// Run executes the root command and handles the application lifecycle.
func Run() {
	// Initialize CLI before running
	InitCLI()

	originalHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		banner.Print()
		originalHelp(c, s)
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
