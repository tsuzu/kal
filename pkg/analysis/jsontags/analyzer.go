package jsontags

import (
	"go/ast"
	"go/types"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// camelCaseRegex is a regular expression that matches camel case strings.
	camelCaseRegex = "^[a-z][a-z0-9]*(?:[A-Z][a-z0-9]*)*$"
)

// Analyzer is the analyzer for the jsontags package.
// It checks that all struct fields in an API are tagged with json tags.
var Analyzer = &analysis.Analyzer{
	Name:     "jsontags",
	Doc:      "Check that all struct fields in an API are tagged with json tags",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		styp, ok := pass.TypesInfo.Types[n.(*ast.StructType)].Type.(*types.Struct)
		// Type information may be incomplete.
		if !ok {
			return
		}

		for i := 0; i < styp.NumFields(); i++ {
			field := styp.Field(i)
			tag := styp.Tag(i)

			checkField(pass, field, tag)
		}
	})

	return nil, nil
}

func checkField(pass *analysis.Pass, field *types.Var, tag string) {
	tagValue, ok := reflect.StructTag(tag).Lookup("json")
	if !ok {
		pass.Reportf(field.Pos(), "field %s is missing json tag", field.Name())
		return
	}

	if tagValue == "" {
		pass.Reportf(field.Pos(), "field %s has empty json tag", field.Name())
		return
	}

	tagValues := strings.Split(tagValue, ",")

	if len(tagValues) == 2 && tagValues[0] == "" && tagValues[1] == "inline" {
		return
	}

	tagName := tagValues[0]
	matched, err := regexp.Match(camelCaseRegex, []byte(tagName))
	if err != nil {
		pass.Reportf(field.Pos(), "error matching json tag: %v", err)
		return
	}

	if !matched {
		pass.Reportf(field.Pos(), "field %s has non-camel case json tag: %s", field.Name(), tagValue)
	}
}
