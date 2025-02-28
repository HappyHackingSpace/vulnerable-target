package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/happyhackingspace/vulnerable-target/internal/config"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/moby/term"
	"github.com/rs/zerolog/log"
)

func Run() {
	settings := config.GetSettings()
	template := templates.Templates[settings.TemplateID]
	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal().Msgf("create client: %v", err)
	}
	defer apiClient.Close()

	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for hostPort, containerPort := range template.Providers["docker"].Ports {
		port, err := nat.NewPort(strings.Split(hostPort, "/")[1], strings.Split(hostPort, "/")[0])
		if err != nil {
			log.Warn().Msgf("invalid port: %v", err)
			continue
		}
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: containerPort}}
	}

	content, cut := strings.CutPrefix(template.Providers["docker"].Content, "IMAGE:")
	if cut {
		imageName := strings.TrimSpace(content)
		_, err := apiClient.ImageInspect(ctx, imageName)
		if err != nil {
			log.Warn().Msgf("inspect image: %v", err)
			pw := progress.NewWriter()
			pw.SetStyle(progress.StyleDefault)
			pw.SetUpdateFrequency(time.Millisecond * 100)
			pw.SetTrackerLength(20)
			pw.Style().Visibility.ETA = false
			pw.Style().Visibility.Percentage = false
			pw.Style().Visibility.Value = false
			go pw.Render()
			tracker := &progress.Tracker{
				Message: fmt.Sprintf("Pulling %s", imageName),
				Total:   0,
				Units:   progress.UnitsDefault,
			}
			pw.AppendTracker(tracker)
			tracker.Start()
			reader, err := apiClient.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				tracker.MarkAsErrored()
				pw.Stop()
				log.Fatal().Msgf("failed to pull image %s: %v", imageName, err)
			}
			defer reader.Close()
			_, err = io.Copy(io.Discard, reader)
			if err != nil {
				log.Fatal().Msgf("error while pulling image %s: %v", imageName, err)
			}
			tracker.MarkAsDone()
			pw.Stop()
		}
		err = createContainer(imageName, apiClient, ctx, exposedPorts, portBindings)
		if err != nil {
			log.Fatal().Msgf("create container: %v", err)
		}
	} else {
		dockerfilePath, err := createDockerfile(content)
		if err != nil {
			log.Fatal().Msgf("create dockerfile: %v", err)
		}

		buildContextTar, err := createBuildContext(dockerfilePath)
		if err != nil {
			log.Fatal().Msgf("create build context: %v", err)
		}
		imageName := fmt.Sprintf("vt-image-%s", template.ID)

		response, err := apiClient.ImageBuild(ctx,
			buildContextTar,
			types.ImageBuildOptions{
				Dockerfile: "Dockerfile",
				Remove:     true,
				Tags:       []string{imageName},
			})
		if err != nil {
			log.Fatal().Msgf("image build: %v", err)
		}
		defer response.Body.Close()

		termFd, isTerm := term.GetFdInfo(os.Stdout)
		err = jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stdout, termFd, isTerm, nil)
		if err != nil {
			log.Fatal().Msgf("display build logs: %v", err)
		}

		err = createContainer(imageName, apiClient, ctx, exposedPorts, portBindings)
		if err != nil {
			log.Fatal().Msgf("create container: %v", err)
		}
	}
}

func createContainer(image string, apiClient *client.Client, ctx context.Context, exposedPorts nat.PortSet, portBindings nat.PortMap) error {
	containerCreate, err := apiClient.ContainerCreate(ctx, &container.Config{
		Image:        image,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return err
	}

	for _, warning := range containerCreate.Warnings {
		log.Warn().Msg(warning)
	}

	err = apiClient.ContainerStart(ctx, containerCreate.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func createDockerfile(content string) (string, error) {
	dir := filepath.Join(os.TempDir(), "vt-dockerfile")

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return "", err
	}

	dockerfilePath := filepath.Join(dir, "Dockerfile")

	err = os.WriteFile(dockerfilePath, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return dockerfilePath, nil
}

func createBuildContext(dockerfilePath string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	dockerfileContent, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, err
	}

	header := &tar.Header{
		Name: "Dockerfile",
		Mode: 0644,
		Size: int64(len(dockerfileContent)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}

	if _, err := tw.Write(dockerfileContent); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}
