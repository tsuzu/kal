package jsontags

import (
	"go/ast"
	"go/types"
	"regexp"

	"github.com/JoelSpeed/kal/pkg/analysis/helpers/extractjsontags"
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
	Requires: []*analysis.Analyzer{inspect.Analyzer, extractjsontags.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	jsonTags := pass.ResultOf[extractjsontags.Analyzer].(extractjsontags.StructFieldTags)

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		s := n.(*ast.StructType)
		styp, ok := pass.TypesInfo.Types[s].Type.(*types.Struct)
		// Type information may be incomplete.
		if !ok {
			return
		}

		for i := 0; i < styp.NumFields(); i++ {
			field := styp.Field(i)

			checkField(pass, s, field, jsonTags)
		}
	})

	return nil, nil
}

func checkField(pass *analysis.Pass, sTyp *ast.StructType, field *types.Var, jsonTags extractjsontags.StructFieldTags) {
	tagInfo := jsonTags.FieldTags(sTyp, field.Name())

	if tagInfo.Missing {
		pass.Reportf(field.Pos(), "field %s is missing json tag", field.Name())
		return
	}

	if tagInfo.Inline {
		return
	}

	if tagInfo.Name == "" {
		pass.Reportf(field.Pos(), "field %s has empty json tag", field.Name())
		return
	}

	matched, err := regexp.Match(camelCaseRegex, []byte(tagInfo.Name))
	if err != nil {
		pass.Reportf(field.Pos(), "error matching json tag: %v", err)
		return
	}

	if !matched {
		pass.Reportf(field.Pos(), "field %s has non-camel case json tag: %s", field.Name(), tagInfo.Name)
	}
}
