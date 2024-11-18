package main

import (
	"github.com/JoelSpeed/kal/pkg/analysis/commentstart"
	"github.com/JoelSpeed/kal/pkg/analysis/jsontags"
	"github.com/JoelSpeed/kal/pkg/analysis/optionalorrequired"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		jsontags.Analyzer,
		commentstart.Analyzer,
		optionalorrequired.Analyzer,
	)
}
