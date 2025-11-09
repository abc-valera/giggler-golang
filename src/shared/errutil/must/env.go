package must

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"time"
)

var ErrInvalidEnvValue = errors.New("invalid env value")

// GetEnv is a shortcut for trimming and empty-cheking environemnt variables.
// Panics if the env var is not set.
func GetEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " env var is not set")
	}

	return strings.TrimSpace(env)
}

// GetEnvBool gets the env var and parses it as bool.
// Panics if the env var is not set or has invalid value.
//
// Possible values:
//   - true
//   - false
func GetEnvBool(key string) bool {
	switch strings.ToLower(GetEnv(key)) {
	case "true":
		return true
	case "false":
		return false
	default:
		panic("failed to parse " + key + " env var as bool")
	}
}

// GetEnvDuration gets the env var and parses it as duration.
// Panics if the env var is not set or has invalid value.
func GetEnvDuration(key string) time.Duration {
	dur, err := time.ParseDuration(GetEnv(key))
	if err != nil {
		panic("failed to parse " + key + " env var as duration")
	}
	return dur
}

// GetEnvUrl gets the env var and parses it as URL.
// Panics if the env var is not set or has invalid value.
func GetEnvUrl(key string) *url.URL {
	return UrlParse(GetEnv(key))
}
