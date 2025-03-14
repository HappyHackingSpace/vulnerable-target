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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/happyhackingspace/vulnerable-target/internal/config"
	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/moby/term"
	"github.com/rs/zerolog/log"
)

func Run() {
	settings := config.GetSettings()
	template := templates.Templates[settings.TemplateID]
	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	defer apiClient.Close()
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for hostPort, containerPort := range template.Providers["docker"].Ports {
		port, err := nat.NewPort(strings.Split(hostPort, "/")[1], strings.Split(hostPort, "/")[0])
		if err != nil {
			log.Warn().Msgf("Invalid port: %v", err)
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
			log.Warn().Msgf("%v", err)

			reader, err := apiClient.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				log.Fatal().Msgf("Failed to pull image %s: %v", imageName, err)
			}
			defer reader.Close()

			termFd, isTerm := term.GetFdInfo(os.Stdout)
			err = jsonmessage.DisplayJSONMessagesStream(reader, os.Stdout, termFd, isTerm, nil)
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
		}
		containerID, err := createContainer(ctx, imageName, apiClient, exposedPorts, portBindings)
		if err != nil {
			if containerID != "" {
				rmErr := apiClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
				if rmErr != nil {
					log.Warn().Msgf("Failed to remove container (%s): %v", containerID, rmErr)
				} else {
					log.Warn().Msgf("Container (%s) is deleted", containerID)
				}
			}
			log.Fatal().Msgf("%v", err)
		}
	} else {
		dockerfilePath, err := createDockerfile(content)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		buildContextTar, err := createBuildContext(dockerfilePath)
		if err != nil {
			log.Fatal().Msgf("%v", err)
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
			log.Fatal().Msgf("%v", err)
		}
		defer response.Body.Close()

		termFd, isTerm := term.GetFdInfo(os.Stdout)
		err = jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stdout, termFd, isTerm, nil)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		containerID, err := createContainer(ctx, imageName, apiClient, exposedPorts, portBindings)
		if err != nil {
			if containerID != "" {
				rmErr := apiClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
				if rmErr != nil {
					log.Warn().Msgf("Failed to remove container (%s): %v", containerID, rmErr)
				} else {
					log.Warn().Msgf("Container (%s) is deleted", containerID)
				}
			}
			log.Fatal().Msgf("%v", err)
		}
	}
}

func createContainer(ctx context.Context, image string, apiClient *client.Client, exposedPorts nat.PortSet, portBindings nat.PortMap) (string, error) {
	containerCreate, err := apiClient.ContainerCreate(ctx, &container.Config{
		Image:        image,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return "", err
	}
	for _, warning := range containerCreate.Warnings {
		log.Warn().Msg(warning)
	}

	err = apiClient.ContainerStart(ctx, containerCreate.ID, container.StartOptions{})
	if err != nil {
		return containerCreate.ID, err
	}

	return containerCreate.ID, nil
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
