package service

import (
	"fmt"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/utils"
)

type action string

// Following some conventions for docker,
// Up is to create container resources
// Down is to remove container resources
const (
	Up   action = "UP"
	Down action = "DOWN"
)

func makeDockerDeploymentJob(config config.Config, action action) (job.Job, error) {
	switch action {
	case Up:
		return createContainersJob(config), nil
	case Down:
		return deleteContainersJob(config), nil
	default:
		return job.Job{}, fmt.Errorf("unknown action %s", action)
	}
}

func createContainersJob(config config.Config) job.Job {
	retryPolicy := utils.GetUIntEnv("DEPLOYMENT_RETRY_POLICY")

	currContainers := make([]container.Kontainer, 0)
	newContainers := make([]container.Kontainer, 0)

	jobsArgs := job.Args{
		"config":         config,
		"currContainers": &currContainers,
		"newContainers":  &newContainers,
	}

	return job.Job{
		ID:          utils.MakeIdentifier(),
		Namespace:   config.Name,
		Type:        ContainerCreate,
		Args:        jobsArgs,
		RetryPolicy: retryPolicy,
		Run:         createContainerResources,
	}
}

func deleteContainersJob(config config.Config) job.Job {
	retryPolicy := utils.GetUIntEnv("DEPLOYMENT_RETRY_POLICY")

	currContainers := make([]container.Kontainer, 0)

	jobsArgs := job.Args{
		"config":         config,
		"currContainers": &currContainers,
	}
	return job.Job{
		ID:          utils.MakeIdentifier(),
		Namespace:   config.Name,
		Type:        ContainerDelete,
		Args:        jobsArgs,
		RetryPolicy: retryPolicy,
		Run:         deleteContainerResources,
	}
}
