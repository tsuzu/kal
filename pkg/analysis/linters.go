package analysis

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

// EnabledByDefaultLinters returns a set of all known linters that are
// enabled by default.
func EnabledByDefaultLinters() sets.Set[string] {
	return sets.New(
		"commentstart",
		"jsontags",
		"optionalorrequired",
	)
}

// DisabledByDefaultLinters returns a set of all known linters that are
// disabled by default.
func DisabledByDefaultLinters() sets.Set[string] {
	return sets.New[string]()
}

// AllLinters returns a set of all known linters, regardless of being
// enabled or disabled by default.
func AllLinters() sets.Set[string] {
	return EnabledByDefaultLinters().Union(DisabledByDefaultLinters())
}
