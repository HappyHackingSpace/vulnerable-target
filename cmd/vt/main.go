package main

import (
	"github.com/happyhackingspace/vulnerable-target/internal/cli"
	"github.com/happyhackingspace/vulnerable-target/internal/logger"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
)

func init() {
	logger.Init()
	templates.Init()
}

func main() {
	cli.Run()
}
