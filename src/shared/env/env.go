package env

import (
	"errors"
	"os"
	"strings"
	"time"
)

var ErrInvalidEnvValue = errors.New("invalid env value")

// Load is a shortcut for trimming and empty-cheking environemnt variables.
// Panics if the env var is not set.
func Load(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " env var is not set")
	}

	return strings.TrimSpace(env)
}

// LoadBool calls Load and casts the value to the boolean type.
//
// Possible values:
//   - true
//   - false
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

// LoadDuration calls Load and parses the duration.
func LoadDuration(key string) time.Duration {
	dur, err := time.ParseDuration(Load(key))
	if err != nil {
		panic("failed to parse " + key + " env var as duration")
	}
	return dur
}
