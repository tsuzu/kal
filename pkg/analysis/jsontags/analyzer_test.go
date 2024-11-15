package jsontags_test

import (
	"testing"

	"github.com/JoelSpeed/kal/pkg/analysis/jsontags"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, jsontags.Analyzer, "a")
}
