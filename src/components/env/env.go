package env

import (
	"errors"
	"os"
	"strings"
	"time"
)

var ErrInvalidEnvValue = errors.New("invalid env value")

// Load is a shortcut for trimming and empty-cheking env vars.
// If the env var is not set, it will exit.
func Load(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " env var is not set")
	}

	return strings.TrimSpace(env)
}

func LoadBool(key string) bool {
	switch strings.ToLower(Load(key)) {
	case "true":
		return true
	case "false":
		return false
	default:
		panic("failed to parse " + key + " env var as bool")
	}
}

func LoadDuration(key string) time.Duration {
	dur, err := time.ParseDuration(Load(key))
	if err != nil {
		panic("failed to parse " + key + " env var as duration")
	}
	return dur
}
