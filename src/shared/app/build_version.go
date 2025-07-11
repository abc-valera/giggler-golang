package app

// buildVersion should be inserted here during the building of the app
// with the use of build tags.
var buildVersion = "unspecified"

func BuildVersion() string {
	return buildVersion
}
