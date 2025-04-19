package main

import (
	"fmt"

	"github.com/happyhackingspace/vulnerable-target/internal/cli"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/pkg/options"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider/dockercompose"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/rs/zerolog/log"
)

func init() {
	logger.Init()
	templates.Init()
	options.LoadEnv()
}

func main() {
	cli.Execute()
	options := options.GetOptions()
	switch options.ProviderName {
	case "docker":
		fmt.Println("docker")
	case "docker-compose":
		if err := (&dockercompose.DockerCompose{}).Start(); err != nil {
			log.Fatal().Msgf("%v", err)
		}
	}
	log.Info().Msgf("%s template is running on %s", options.TemplateID, options.ProviderName)
	fmt.Println("Hello, World!")
}
