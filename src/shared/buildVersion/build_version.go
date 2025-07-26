package buildVersion

// buildVersion should be inserted here during the building of the app
// with the use of build tags.
var buildVersion = "unspecified"

func Get() string {
	return buildVersion
}
