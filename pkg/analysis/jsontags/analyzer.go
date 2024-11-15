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
	Name:       "jsontags",
	Doc:        "Check that all struct fields in an API are tagged with json tags",
	Run:        run,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	ResultType: reflect.TypeOf(new(map[ast.Node]map[string]FieldTagInfo)),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	results := make(map[ast.Node]map[string]FieldTagInfo)

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		results[n] = make(map[string]FieldTagInfo)

		styp, ok := pass.TypesInfo.Types[n.(*ast.StructType)].Type.(*types.Struct)
		// Type information may be incomplete.
		if !ok {
			return
		}

		for i := 0; i < styp.NumFields(); i++ {
			field := styp.Field(i)
			tag := styp.Tag(i)

			results[n][field.Name()] = checkField(pass, field, tag)
		}
	})

	return &results, nil
}

func checkField(pass *analysis.Pass, field *types.Var, tag string) FieldTagInfo {
	tagValue, ok := reflect.StructTag(tag).Lookup("json")
	if !ok {
		pass.Reportf(field.Pos(), "field %s is missing json tag", field.Name())
		return FieldTagInfo{}
	}

	if tagValue == "" {
		pass.Reportf(field.Pos(), "field %s has empty json tag", field.Name())
		return FieldTagInfo{}
	}

	tagValues := strings.Split(tagValue, ",")

	if len(tagValues) == 2 && tagValues[0] == "" && tagValues[1] == "inline" {
		return FieldTagInfo{Inline: true}
	}

	tagName := tagValues[0]
	matched, err := regexp.Match(camelCaseRegex, []byte(tagName))
	if err != nil {
		pass.Reportf(field.Pos(), "error matching json tag: %v", err)
		return FieldTagInfo{}
	}

	if !matched {
		pass.Reportf(field.Pos(), "field %s has non-camel case json tag: %s", field.Name(), tagValue)
	}

	return FieldTagInfo{Name: tagName, OmitEmpty: len(tagValues) == 2 && tagValues[1] == "omitempty"}
}

// FieldTagInfo contains information about a field's json tag.
// This is used to pass information about a field's json tag between analyzers.
type FieldTagInfo struct {
	// Name is the name of the field extracted from the json tag.
	Name string

	// OmitEmpty is true if the field has the omitempty option in the json tag.
	OmitEmpty bool

	// Inline is true if the field has the inline option in the json tag.
	Inline bool
}
