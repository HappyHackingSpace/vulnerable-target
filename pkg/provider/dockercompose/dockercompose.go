// Package dockercompose provides Docker Compose provider implementation for managing vulnerable target environments.
package dockercompose

import (
	"fmt"

	"github.com/happyhackingspace/vulnerable-target/internal/state"
	"github.com/happyhackingspace/vulnerable-target/pkg/provider"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/disk"
	tmpl "github.com/happyhackingspace/vulnerable-target/pkg/template"
)

var _ provider.Provider = &DockerCompose{}

// DockerCompose implements the Provider interface using Docker Compose.
type DockerCompose struct{}

// Name returns the provider name.
func (d *DockerCompose) Name() string {
	return "docker-compose"
}

// Start launches the vulnerable target environment using Docker Compose.
func (d *DockerCompose) Start(template *tmpl.Template) error {
	cfg := disk.NewConfig().WithFileName("deployments.db").WithBucketName("deployments")
	st, err := state.NewManager(cfg)
	if err != nil {
		return err
	}

	exist, _ := st.DeploymentExist(d.Name(), template.ID) //nolint:errcheck
	if exist {
		return fmt.Errorf("already running")
	}

	dockerCli, err := createDockerCLI()
	if err != nil {
		return err
	}

	project, err := loadComposeProject(*template)
	if err != nil {
		return err
	}

	err = runComposeUp(dockerCli, project)
	if err != nil {
		return err
	}

	err = st.AddNewDeployment(d.Name(), template.ID)
	if err != nil {
		return err
	}

	return nil
}

// Stop shuts down the vulnerable target environment using Docker Compose.
func (d *DockerCompose) Stop(template *tmpl.Template) error {
	cfg := disk.NewConfig().WithFileName("deployments.db").WithBucketName("deployments")
	st, err := state.NewManager(cfg)
	if err != nil {
		return err
	}

	exist, err := st.DeploymentExist(d.Name(), template.ID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("deployment not exist")
	}

	dockerCli, err := createDockerCLI()
	if err != nil {
		return err
	}

	project, err := loadComposeProject(*template)
	if err != nil {
		return err
	}

	err = runComposeDown(dockerCli, project)
	if err != nil {
		return err
	}

	err = st.RemoveDeployment(d.Name(), template.ID)
	if err != nil {
		return err
	}

	return nil
}

// Status returns status the vulnerable target environment using Docker Compose.
func (d *DockerCompose) Status(template *tmpl.Template) (string, error) {
	dockerCli, err := createDockerCLI()
	if err != nil {
		return "unknown", err
	}

	project, err := loadComposeProject(*template)
	if err != nil {
		return "unknown", err
	}

	running, err := runComposeStats(dockerCli, project)
	if err != nil {
		return "unknown", err
	}

	if !running {
		return "unknown", err
	}

	return "running", err
}
