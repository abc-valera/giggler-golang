package appMode

import (
	"github.com/abc-valera/giggler-golang/src/components/env"
)

var appModeVar string

func init() {
	switch appModeEnv := env.Load("APP_MODE"); appModeEnv {
	case "dev", "prod", "test":
		appModeVar = appModeEnv
	default:
		panic("APP_MODE environment variable is invalid")
	}
}

func IsDev() bool {
	return appModeVar == "dev"
}

func IsProd() bool {
	return appModeVar == "prod"
}

func IsTest() bool {
	return appModeVar == "test"
}
