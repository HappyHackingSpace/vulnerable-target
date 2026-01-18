package dockercompose

import (
	"context"
	"fmt"
	"time"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/happyhackingspace/vt/internal/app"
	tmpl "github.com/happyhackingspace/vt/pkg/template"
)

func createDockerCLI() (command.Cli, error) {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}

	opts := flags.NewClientOptions()
	err = dockerCli.Initialize(opts)
	if err != nil {
		return nil, err
	}

	return dockerCli, nil
}

func loadComposeProject(template tmpl.Template) (*types.Project, error) {
	cfg := app.DefaultConfig()
	composePath, workingDir, err := tmpl.GetDockerComposePath(template.ID, cfg.TemplatesPath)
	if err != nil {
		return nil, err
	}

	projectName := fmt.Sprintf("vt-compose-%s", template.ID)

	configDetails := types.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []types.ConfigFile{
			{
				Filename: composePath,
			},
		},
		Environment: map[string]string{
			"COMPOSE_PROJECT_NAME": projectName,
		},
	}

	project, err := loader.LoadWithContext(
		context.Background(),
		configDetails,
		func(options *loader.Options) {
			options.SkipValidation = false
			options.SkipInterpolation = false
			options.SetProjectName(projectName, true)
			options.ResolvePaths = true
		},
	)
	if err != nil {
		return nil, err
	}

	updatedServices := make(types.Services, len(project.Services))
	for name, service := range project.Services {
		serviceCopy := service
		if serviceCopy.Labels == nil {
			serviceCopy.Labels = make(map[string]string)
		}
		serviceCopy.Labels["com.docker.compose.project"] = projectName
		serviceCopy.Labels["com.docker.compose.service"] = name
		serviceCopy.Labels["com.docker.compose.project.working_dir"] = workingDir
		serviceCopy.Labels["com.docker.compose.project.config_files"] = composePath
		serviceCopy.Labels["com.docker.compose.config-hash"] = name
		serviceCopy.Labels["com.docker.compose.oneoff"] = "False"
		updatedServices[name] = serviceCopy
	}
	project.Services = updatedServices

	return project, nil
}

func runComposeUp(dockerCli command.Cli, project *types.Project) error {
	composeService := compose.NewComposeService(dockerCli)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err := composeService.Pull(ctx, project, api.PullOptions{})
	if err != nil {
		return err
	}

	err = composeService.Create(ctx, project, api.CreateOptions{
		Services:      project.ServiceNames(),
		RemoveOrphans: true,
		Recreate:      api.RecreateForce,
		Inherit:       true,
		QuietPull:     false,
	})
	if err != nil {
		return err
	}

	err = composeService.Start(ctx, project.Name, api.StartOptions{
		Project:  project,
		Attach:   nil,
		Services: project.ServiceNames(),
	})
	if err != nil {
		return err
	}

	return nil
}

func runComposeDown(dockerCli command.Cli, project *types.Project) error {
	composeService := compose.NewComposeService(dockerCli)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err := composeService.Down(ctx, project.Name, api.DownOptions{
		Project:       project,
		RemoveOrphans: true,
		Volumes:       true,
	})
	if err != nil {
		return err
	}

	return nil
}

func runComposeStats(dockerCli command.Cli, project *types.Project) (bool, error) {
	composeService := compose.NewComposeService(dockerCli)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	summary, err := composeService.Ps(ctx, project.Name, api.PsOptions{
		Project: project,
		All:     true,
	})

	if err != nil {
		return false, err
	}

	for i := 0; i < len(summary); i++ {
		if summary[i].State != "running" {
			return false, nil
		}
	}

	return true, nil
}
