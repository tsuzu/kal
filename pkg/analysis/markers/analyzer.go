package markers

import (
	"go/ast"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the analyzer for the markers package.
// It iterates over declarations within a package and parses the comments to extract markers.
var Analyzer = &analysis.Analyzer{
	Name:       "markers",
	Doc:        "Iterates over declarations within a package and parses the comments to extract markers",
	Run:        run,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	ResultType: reflect.TypeOf(new(Markers)),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter to declarations so that we can look at all types in the package.
	declFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	results := &Markers{
		StructMarkers: make(map[*ast.StructType]StructMarkers),
	}

	inspect.Preorder(declFilter, func(n ast.Node) {
		decl, ok := n.(*ast.GenDecl)
		if !ok {
			return
		}

		for _, spec := range decl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			switch typeSpec.Type.(type) {
			case *ast.StructType:
				sTyp := typeSpec.Type.(*ast.StructType)
				results.StructMarkers[sTyp] = extractStructMarkers(decl, sTyp)
			}
		}
	})

	return results, nil
}

func extractStructMarkers(decl *ast.GenDecl, sTyp *ast.StructType) StructMarkers {
	structMarkers := StructMarkers{}

	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if marker := extractMarker(comment.Text); marker != "" {
				structMarkers.Markers = append(structMarkers.Markers, marker)
			}
		}
	}

	if sTyp.Fields == nil {
		return structMarkers
	}

	structMarkers.FieldMarkers = make(map[string][]string)

	for _, field := range sTyp.Fields.List {
		if field == nil || len(field.Names) == 0 {
			continue
		}

		if field.Doc == nil {
			continue
		}

		fieldName := field.Names[0].Name
		for _, comment := range field.Doc.List {
			if marker := extractMarker(comment.Text); marker != "" {
				structMarkers.FieldMarkers[fieldName] = append(structMarkers.FieldMarkers[fieldName], marker)
			}
		}
	}

	return structMarkers
}

func extractMarker(comment string) string {
	if strings.HasPrefix(comment, "// +") {
		return strings.TrimPrefix(comment, "// +")
	}

	return ""
}

type Markers struct {
	// StructMarkers contains the markers for each struct in the package.
	StructMarkers map[*ast.StructType]StructMarkers
}

type StructMarkers struct {
	// Markers contains the markers for the struct.
	Markers []string

	// FieldMarkers contains the markers for each field in the struct.
	// Mapped by field name.
	FieldMarkers map[string][]string
}
