package env

import (
	"errors"
	"os"
	"strings"
	"time"
)

var ErrInvalidEnvValue = errors.New("invalid env value")

// Load is a shortcut for trimming and empty-cheking environemnt variables.
// If the environment variable is not set, it will exit.
func Load(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " environment variable is not set")
	}

	return strings.TrimSpace(env)
}

func LoadDuration(key string) time.Duration {
	dur, err := time.ParseDuration(Load(key))
	if err != nil {
		panic("failed to parse " + key + " environment variable as duration")
	}
	return dur
}
