package config

// LintersConfig contains configuration for individual linters.
type LintersConfig struct {
	// jsonTags contains configuration for the jsontags linter.
	JSONTags JSONTagsConfig `json:"jsonTags"`
}

// JSONTagsConfig contains configuration for the jsontags linter.
type JSONTagsConfig struct {
	// jsonTagRegex is the regular expression used to validate that json tags are in a particular format.
	// By default, the regex used is "^[a-z][a-z0-9]*(?:[A-Z][a-z0-9]*)*$" and is used to check for
	// camel case like string.
	JSONTagRegex string `json:"jsonTagRegex"`
}
