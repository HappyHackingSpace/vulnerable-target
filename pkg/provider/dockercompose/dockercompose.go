package dockercompose

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/happyhackingspace/vulnerable-target/internal/file"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
)

var _ provider.Provider = &DockerCompose{}

type DockerCompose struct{}

func (d *DockerCompose) Name() string {
	return "docker-compose"
}

func (d *DockerCompose) Start() error {
	template, err := templates.GetCurrentTemplate()
	if err != nil {
		return err
	}

	composeContent := template.Providers["docker_compose"].Content

	composeFilePath, err := file.CreateTempFile(composeContent, "docker-compose.yml")
	if err != nil {
		return err
	}

	upCmd := exec.Command("docker", "compose", "-f", composeFilePath, "-p", fmt.Sprintf("vt-compose-%s", template.ID), "up", "-d")
	upCmd.Stdout = os.Stdout
	upCmd.Stderr = os.Stderr

	err = upCmd.Run()
	if err != nil {
		return err
	}

	err = file.DeteleFile(composeFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerCompose) Stop() error {
	template, err := templates.GetCurrentTemplate()
	if err != nil {
		return err
	}

	composeContent := template.Providers["docker_compose"].Content

	composeFilePath, err := file.CreateTempFile(composeContent, "docker-compose.yml")
	if err != nil {
		return err
	}

	downCmd := exec.Command("docker", "compose", "-f", composeFilePath, "-p", fmt.Sprintf("vt-compose-%s", template.ID), "down", "--volumes")
	downCmd.Stdout = os.Stdout
	downCmd.Stderr = os.Stderr

	err = downCmd.Run()
	if err != nil {
		return err
	}

	err = file.DeteleFile(composeFilePath)
	if err != nil {
		return err
	}

	return nil
}
