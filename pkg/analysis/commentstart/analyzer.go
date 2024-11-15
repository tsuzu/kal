package commentstart

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/JoelSpeed/kal/pkg/analysis/jsontags"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the analyzer for the commentstart package.
// It checks that all struct fields in an API have a godoc, and that the godoc starts with the serialised field name.
var Analyzer = &analysis.Analyzer{
	Name:     "commentstart",
	Doc:      "Check that all struct fields in an API have a godoc, and that the godoc starts with the serialised field name",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, jsontags.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	jsonTags := pass.ResultOf[jsontags.Analyzer].(*map[ast.Node]map[string]jsontags.FieldTagInfo)

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		jsonTagsByField := (*jsonTags)[n]

		styp, ok := n.(*ast.StructType)
		if !ok {
			return
		}

		if styp.Fields == nil {
			return
		}

		for _, field := range styp.Fields.List {
			checkField(pass, field, jsonTagsByField)
		}
	})

	return nil, nil
}

func checkField(pass *analysis.Pass, field *ast.Field, jsonTagsByField map[string]jsontags.FieldTagInfo) {
	if field == nil || len(field.Names) == 0 {
		return
	}

	fieldName := field.Names[0].Name
	tagInfo := jsonTagsByField[fieldName]

	if tagInfo.Name == "" {
		return
	}

	if field.Doc == nil {
		pass.Reportf(field.Pos(), "field %s is missing godoc comment", fieldName)
		return
	}

	firstLine := field.Doc.List[0]
	if !strings.HasPrefix(firstLine.Text, "// "+tagInfo.Name+" ") {
		if strings.HasPrefix(strings.ToLower(firstLine.Text), strings.ToLower("// "+tagInfo.Name+" ")) {
			// The comment start is correct, apart from the casing, we can fix that.
			pass.Report(analysis.Diagnostic{
				Pos:     firstLine.Pos(),
				Message: fmt.Sprintf("godoc for field %s should start with '%s ...'", fieldName, tagInfo.Name),
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: fmt.Sprintf("should replace first word with `%s`", tagInfo.Name),
						TextEdits: []analysis.TextEdit{
							{
								Pos:     firstLine.Pos(),
								End:     firstLine.Pos() + token.Pos(len(tagInfo.Name)) + token.Pos(4),
								NewText: []byte("// " + tagInfo.Name + " "),
							},
						},
					},
				},
			})
		} else {
			pass.Reportf(field.Doc.List[0].Pos(), "godoc for field %s should start with '%s ...'", fieldName, tagInfo.Name)
		}
	}
}
