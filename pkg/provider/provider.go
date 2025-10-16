// Package provider defines interfaces and types for managing vulnerable target environments.
package provider

import (
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
)

// Provider defines the interface for managing vulnerable target environments.
type Provider interface {
	Name() string
	Start(template *templates.Template) error
	Stop(template *templates.Template) error
	Status(template *templates.Template) (string, error)
}
