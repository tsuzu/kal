package validation

import (
	"fmt"
	"regexp"

	"github.com/JoelSpeed/kal/pkg/config"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateLintersConfig is used to validate the configuration in the config.LintersConfig struct.
func ValidateLintersConfig(lc config.LintersConfig, fldPath *field.Path) field.ErrorList {
	fieldErrors := field.ErrorList{}

	fieldErrors = append(fieldErrors, validateJSONTagsConfig(lc.JSONTags, fldPath.Child("jsonTags"))...)

	return fieldErrors
}

// validateJSONTagsConfig is used to validate the configuration in the config.JSONTagsConfig struct.
func validateJSONTagsConfig(jtc config.JSONTagsConfig, fldPath *field.Path) field.ErrorList {
	fieldErrors := field.ErrorList{}

	if jtc.JSONTagRegex != "" {
		if _, err := regexp.Compile(jtc.JSONTagRegex); err != nil {
			fieldErrors = append(fieldErrors, field.Invalid(fldPath.Child("jsonTagRegex"), jtc.JSONTagRegex, fmt.Sprintf("invalid regex: %v", err)))
		}
	}

	return fieldErrors
}
