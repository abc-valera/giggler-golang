package must

import "net/url"

// ParseUrl tries to parse the provided URL and panics on error.
func ParseUrl(rawurl string) *url.URL {
	return Do(url.Parse(rawurl))
}
