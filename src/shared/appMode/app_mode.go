package appMode

import (
	"github.com/abc-valera/giggler-golang/src/components/env"
)

var appModeVar = func() string {
	switch appModeEnv := env.Load("APP_MODE"); appModeEnv {
	case "dev", "prod", "test":
		return appModeEnv
	default:
		panic("APP_MODE env var is invalid")
	}
}()

func IsDev() bool {
	return appModeVar == "dev"
}

func IsProd() bool {
	return appModeVar == "prod"
}

func IsTest() bool {
	return appModeVar == "test"
}
