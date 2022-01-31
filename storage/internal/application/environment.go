package application

import "os"

type Environment string

const (
	Development Environment = `Development`
	Production  Environment = `Production`

	ServiceName = `storage`
)

func NewEnvironment() Environment {
	name, ok := os.LookupEnv(`ENV`)
	if !ok || !Environment(name).Is(Production) {
		return Development
	}

	return Production
}

func (e Environment) Is(env Environment) bool {
	return e == env
}
