package must

import "net/url"

// UrlParse tries to parse the provided URL and panics on error.
func UrlParse(rawurl string) *url.URL {
	return Do(url.Parse(rawurl))
}
