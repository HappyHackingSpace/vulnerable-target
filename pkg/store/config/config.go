// Package config provides configuration interfaces and base implementations for storage backends
package config

// Interface defines the contract for configuration objects
type Interface interface {
	GetName() string
}

// Struct provides a base implementation of the Interface with a name field
type Struct struct {
	name string
}

// GetName returns the name of the configuration
func (s Struct) GetName() string {
	return s.name
}
