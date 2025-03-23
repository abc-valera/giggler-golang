package app

import (
	"github.com/abc-valera/giggler-golang/src/shared/env"
)

var modeVar = initMode()

func initMode() mode {
	switch modeEnv := mode(env.Load("APP_MODE")); modeEnv {
	case dev, prod, test:
		return modeEnv
	default:
		panic("APP_MODE env var is invalid")
	}
}

func Mode() mode {
	return modeVar
}

type mode string

var (
	dev  mode = "dev"
	prod mode = "prod"
	test mode = "test"
)

func (m mode) IsDev() bool {
	return m == dev
}

func (m mode) IsProd() bool {
	return m == prod
}

func (m mode) IsTest() bool {
	return m == test
}
