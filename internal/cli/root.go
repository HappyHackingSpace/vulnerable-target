package cli

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/internal/options"
	"github.com/happyhackingspace/vulnerable-target/internal/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var LogLevels = map[string]bool{
	zerolog.DebugLevel.String(): true,
	zerolog.InfoLevel.String():  true,
	zerolog.WarnLevel.String():  true,
	zerolog.ErrorLevel.String(): true,
	zerolog.FatalLevel.String(): true,
	zerolog.PanicLevel.String(): true,
}

func init() {
	options := options.GetOptions()

	rootCmd.PersistentFlags().StringVarP(&options.VerbosityLevel, "verbosity", "v", zerolog.InfoLevel.String(),
		fmt.Sprintf("Set the verbosity level for logs (%s)",
			strings.Join(slices.Collect(maps.Keys(LogLevels)), ", ")))
}

var rootCmd = &cobra.Command{
	Use:     "vt",
	Short:   "Create vulnerable environment",
	Version: utils.AppVersion,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		logger.Init()
		if cmd.Name() != "help" {
			fmt.Println(utils.Banner())
		}
	},
	SilenceErrors: true,
}

func Run() {

	originalHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		fmt.Println(utils.Banner())
		originalHelp(c, s)
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
