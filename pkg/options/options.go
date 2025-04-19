package options

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Options struct {
	VerbosityLevel string
	ProviderName   string
	TemplateID     string
}

var GlobalOptions Options

func GetOptions() *Options {
	return &GlobalOptions
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msgf("Error loading .env file: %v", err)
	}
}
