package provider

type Provider interface {
	Name() string
	Start() error
	Stop() error
}
