package optionalorrequired_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/optionalorrequired"
	"github.com/JoelSpeed/kal/pkg/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDefaultConfiguration(t *testing.T) {
	testdata := analysistest.TestData()

	a, err := optionalorrequired.Initializer().Init(config.LintersConfig{})
	if err != nil {
		t.Fatal(err)
	}

	analysistest.RunWithSuggestedFixes(t, testdata, a, "a")
}

func TestSwappedMarkerPriority(t *testing.T) {
	testdata := analysistest.TestData()

	a, err := optionalorrequired.Initializer().Init(config.LintersConfig{
		OptionalOrRequired: config.OptionalOrRequiredConfig{
			PreferredOptionalMarker: optionalorrequired.KubebuilderOptionalMarker,
			PreferredRequiredMarker: optionalorrequired.KubebuilderRequiredMarker,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	analysistest.RunWithSuggestedFixes(t, testdata, a, "b")
}
