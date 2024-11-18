package optionalorrequired_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/optionalorrequired"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, optionalorrequired.Analyzer, "a")
}
