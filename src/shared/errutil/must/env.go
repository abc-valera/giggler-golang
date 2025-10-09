package must

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"time"
)

var ErrInvalidEnvValue = errors.New("invalid env value")

// Env is a shortcut for trimming and empty-cheking environemnt variables.
// Panics if the env var is not set.
func Env(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " env var is not set")
	}

	return strings.TrimSpace(env)
}

// EnvBool gets the env var and parses it as bool.
// Panics if the env var is not set or has invalid value.
//
// Possible values:
//   - true
//   - false
func EnvBool(key string) bool {
	switch strings.ToLower(Env(key)) {
	case "true":
		return true
	case "false":
		return false
	default:
		panic("failed to parse " + key + " env var as bool")
	}
}

// EnvDuration gets the env var and parses it as duration.
// Panics if the env var is not set or has invalid value.
func EnvDuration(key string) time.Duration {
	dur, err := time.ParseDuration(Env(key))
	if err != nil {
		panic("failed to parse " + key + " env var as duration")
	}
	return dur
}

// EnvUrl gets the env var and parses it as URL.
// Panics if the env var is not set or has invalid value.
func EnvUrl(key string) *url.URL {
	return UrlParse(Env(key))
}
