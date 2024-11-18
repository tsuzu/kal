package markers

import (
	"go/ast"
	"go/token"
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
	structMarkers := StructMarkers{
		Markers: NewMarkerSet(),
	}

	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if marker := extractMarker(comment); marker.Value != "" {
				structMarkers.Markers.Insert(marker)
			}
		}
	}

	if sTyp.Fields == nil {
		return structMarkers
	}

	structMarkers.FieldMarkers = make(map[string]MarkerSet)

	for _, field := range sTyp.Fields.List {
		if field == nil || len(field.Names) == 0 {
			continue
		}

		if field.Doc == nil {
			continue
		}

		fieldName := field.Names[0].Name
		for _, comment := range field.Doc.List {
			if marker := extractMarker(comment); marker.Value != "" {
				if structMarkers.FieldMarkers[fieldName] == nil {
					structMarkers.FieldMarkers[fieldName] = NewMarkerSet()
				}

				structMarkers.FieldMarkers[fieldName].Insert(marker)
			}
		}
	}

	return structMarkers
}

func extractMarker(comment *ast.Comment) Marker {
	if !strings.HasPrefix(comment.Text, "// +") {
		return Marker{}
	}

	return Marker{
		Value:      strings.TrimPrefix(comment.Text, "// +"),
		RawComment: comment.Text,
		Pos:        comment.Pos(),
		End:        comment.End(),
	}
}

type Markers struct {
	// StructMarkers contains the markers for each struct in the package.
	StructMarkers map[*ast.StructType]StructMarkers
}

type StructMarkers struct {
	// Markers contains the markers for the struct.
	// It uses the MarkerSet to store the markers allowing lookup of detailed
	// marker information based on the marker value.
	Markers MarkerSet

	// FieldMarkers contains the markers for each field in the struct.
	// Mapped  by field name to a MarkerSet for the field.
	// MarkerSet stores the markers, allowing lookup of detailed
	// marker information based on the marker value.
	FieldMarkers map[string]MarkerSet
}

type Marker struct {
	// Value is the value of the marker once the leading comment and '+' are trimmed.
	Value string

	// RawComment is the raw comment line, unfiltered.
	RawComment string

	// Pos is the starting position in the file for the comment line containing the marker.
	Pos token.Pos

	// End is the ending position in the file for the coment line containing the marker.
	End token.Pos
}

// MarkerSet is a set implementation for Markers that uses
// the Marker value as the key, but returns the full Marker
// as the result.
type MarkerSet map[string]Marker

// NewMarkerSet initialises a new MarkerSet with the provided values.
// If any markers have the same value, the latter marker in the list
// will take precedence, no duplication checks are implemented.
func NewMarkerSet(markers ...Marker) MarkerSet {
	ms := make(MarkerSet)

	ms.Insert(markers...)

	return ms
}

// Insert add the given markers to the MarkerSet.
// If any markers have the same value, the latter marker in the list
// will take precedence, no duplication checks are implemented.
func (ms MarkerSet) Insert(markers ...Marker) {
	for _, marker := range markers {
		ms[marker.Value] = marker
	}
}

// Has returns whether a marker with the value given is present in the
// MarkerSet.
func (ms MarkerSet) Has(value string) bool {
	_, ok := ms[value]
	return ok
}
