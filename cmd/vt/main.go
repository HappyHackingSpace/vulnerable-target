// Package main is the entry point for the vulnerable target application.
package main

import (
	"github.com/happyhackingspace/vulnerable-target/internal/banner"
	"github.com/happyhackingspace/vulnerable-target/internal/cli"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
)

func main() {
	// Initialize logger and templates explicitly
	logger.Init()
	templates.Init()

	// Print the banner
	banner.Print()

	// Run the CLI
	cli.Run()
}
