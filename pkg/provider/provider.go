// Package provider defines interfaces and types for managing vulnerable target environments.
package provider

import (
	tmpl "github.com/happyhackingspace/vulnerable-target/pkg/template"
)

// Provider defines the interface for managing vulnerable target environments.
type Provider interface {
	Name() string
	Start(template *tmpl.Template) error
	Stop(template *tmpl.Template) error
	Status(template *tmpl.Template) (string, error)
}
