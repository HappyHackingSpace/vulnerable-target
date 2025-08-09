package options

type Options struct {
	VerbosityLevel string
	ProviderName   string
	TemplateID     string
}

var GlobalOptions Options

func GetOptions() *Options {
	return &GlobalOptions
}
