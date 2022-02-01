package application

import "os"

const (
	Development Environment = `Development`
	Production  Environment = `Production`

	ServiceName = `wallet`
)

type Environment string

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
