package infra

import "context"

type Container struct {
	Config *Config
}

func NewContainer(config *Config, ctx context.Context) *Container {
	container := &Container{}

	container.Config = config

	return container
}
