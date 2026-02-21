package infra

type Container struct {
	Config *Config
}

func NewContainer(config *Config) *Container {
	container := &Container{}
	container.Config = config

	return container
}
