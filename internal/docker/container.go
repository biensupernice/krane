package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// CreateContainerConfig : properties required to create a container
type CreateContainerConfig struct {
	ContainerName string
	Image         string
	NetworkID     string
	Labels        map[string]string
	Ports         nat.PortMap
	VolumesMount  []mount.Mount
	VolumesMap    map[string]struct{}
	Env           []string // Comma separated, formatted NODE_ENV=dev
	Command       []string
	Entrypoint    []string
}

// CreateContainer : create a docker container
func (c *Client) CreateContainer(
	ctx context.Context,
	config CreateContainerConfig,
) (container.ContainerCreateCreatedBody, error) {
	networkingConfig := createNetworkingConfig(config.NetworkID)
	hostConfig := createHostConfig(config.Ports, config.VolumesMount)
	containerConfig := createContainerConfig(config.ContainerName,
		config.Image,
		config.Env,
		config.Labels,
		config.Command,
		config.Entrypoint,
		config.VolumesMap)

	return c.ContainerCreate(
		ctx,
		&containerConfig,
		&hostConfig,
		&networkingConfig,
		config.ContainerName,
	)
}

// StartContainer : start a docker container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	options := types.ContainerStartOptions{}
	return c.ContainerStart(ctx, containerID, options)
}

// StopContainer : stop docker container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	timeout := 60 * time.Second
	return c.ContainerStop(ctx, containerID, &timeout)
}

// RemoveContainer : remove docker container
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         force,
	}
	return c.ContainerRemove(ctx, containerID, options)
}

// GetOneContainer : get one container
func (c *Client) GetOneContainer(ctx context.Context, containerId string) (types.ContainerJSON, error) {
	return c.ContainerInspect(ctx, containerId)
}

// GetKraneContainers : gets all containers on the host machine
func (c *Client) GetAllContainers(ctx *context.Context) ([]types.ContainerJSON, error) {
	options := types.ContainerListOptions{
		All:   true,
		Quiet: false,
	}

	containers, err := c.ContainerList(*ctx, options)
	if err != nil {
		return make([]types.ContainerJSON, 0), err
	}

	toJsonContainers := make([]types.ContainerJSON, 0)
	for _, container := range containers {
		containerJson, err := c.GetOneContainer(*ctx, container.ID)
		if err != nil {
			return make([]types.ContainerJSON, 0), err
		}

		toJsonContainers = append(toJsonContainers, containerJson)
	}
	return toJsonContainers, nil
}

// GetContainerStatus : get the response of a container
func (c *Client) GetContainerStatus(ctx context.Context, containerID string, stream bool) (stats types.ContainerStats, err error) {
	return c.ContainerStats(ctx, containerID, stream)
}

// StreamContainerLogs : stream container logs
func (c *Client) StreamContainerLogs(ctx *context.Context, containerID string) (reader io.Reader, err error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "50",
	}

	return c.ContainerLogs(*ctx, containerID, options)
}

// ConnectContainerToNetwork : connect a container to a network
func (c *Client) ConnectContainerToNetwork(ctx *context.Context, networkID string, containerID string) (err error) {
	config := network.EndpointSettings{NetworkID: networkID}
	return c.NetworkConnect(*ctx, networkID, containerID, &config)
}

func createContainerConfig(
	hostname string,
	image string,
	env []string,
	labels map[string]string,
	command []string,
	entrypoint []string,
	volumes map[string]struct{}) container.Config {
	config := container.Config{
		Hostname: hostname,
		Image:    image,
		Env:      env,
		Labels:   labels,
		Volumes:  volumes,
	}

	if len(command) > 0 {
		config.Cmd = command
	}

	if len(entrypoint) > 0 {
		config.Entrypoint = entrypoint
	}

	return config
}

func createHostConfig(ports nat.PortMap, volumes []mount.Mount) container.HostConfig {
	return container.HostConfig{
		PortBindings: ports,
		AutoRemove:   true,
		Mounts:       volumes,
	}
}
