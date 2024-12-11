package markers

import (
	"errors"
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	errCouldNotGetInspector  = errors.New("could not get inspector")
	errCouldNotCreateMarkers = errors.New("could not create markers")
)

// Markers allows access to markers extracted from the
// go types.
type Markers interface {
	// StructMarkers returns markers associated to the given sturct.
	StructMarkers(*ast.StructType) MarkerSet

	// StructFieldMarkers returns markers associated to the named field in the given struct.
	StructFieldMarkers(*ast.StructType, string) MarkerSet
}

func newMarkers() Markers {
	return &markers{
		structMarkers:      make(map[*ast.StructType]MarkerSet),
		structFieldMarkers: make(map[*ast.StructType]map[string]MarkerSet),
	}
}

// markers implements the storage for the implementation of the Markers interface.
type markers struct {
	structMarkers      map[*ast.StructType]MarkerSet
	structFieldMarkers map[*ast.StructType]map[string]MarkerSet
}

// StructMarkers returns the appropriate MarkerSet if found, else
// it returns an empty MarkerSet.
func (m *markers) StructMarkers(sTyp *ast.StructType) MarkerSet {
	sMarkers, ok := m.structMarkers[sTyp]
	if !ok {
		return NewMarkerSet()
	}

	return sMarkers
}

// StructFieldMarkers return the appropriate MarkerSet for the named field in the
// given struct, or an empty MarkerSet if the appropriate MarkerSet isn't found.
func (m *markers) StructFieldMarkers(sTyp *ast.StructType, field string) MarkerSet {
	sMarkers, ok := m.structFieldMarkers[sTyp]
	if !ok {
		return NewMarkerSet()
	}

	fMarkers, ok := sMarkers[field]
	if !ok {
		return NewMarkerSet()
	}

	return fMarkers
}

func (m *markers) insertStructMarkers(sTyp *ast.StructType, ms MarkerSet) {
	m.structMarkers[sTyp] = ms
}

func (m *markers) insertStructFieldMarkers(sTyp *ast.StructType, field string, ms MarkerSet) {
	if m.structFieldMarkers[sTyp] == nil {
		m.structFieldMarkers[sTyp] = make(map[string]MarkerSet)
	}

	m.structFieldMarkers[sTyp][field] = ms
}

// Analyzer is the analyzer for the markers package.
// It iterates over declarations within a package and parses the comments to extract markers.
var Analyzer = &analysis.Analyzer{
	Name:       "markers",
	Doc:        "Iterates over declarations within a package and parses the comments to extract markers",
	Run:        run,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	ResultType: reflect.TypeOf(newMarkers()),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errCouldNotGetInspector
	}

	// Filter to declarations so that we can look at all types in the package.
	declFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	results, ok := newMarkers().(*markers)
	if !ok {
		return nil, errCouldNotCreateMarkers
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

			// TODO: Add support for other types and remove nolint.
			switch typeSpec.Type.(type) { //nolint:gocritic
			case *ast.StructType:
				sTyp, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				extractStructMarkers(decl, sTyp, results)
			}
		}
	})

	return results, nil
}

func extractStructMarkers(decl *ast.GenDecl, sTyp *ast.StructType, results *markers) {
	structMarkers := NewMarkerSet()

	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if marker := extractMarker(comment); marker.Value != "" {
				structMarkers.Insert(marker)
			}
		}
	}

	results.insertStructMarkers(sTyp, structMarkers)

	for _, field := range sTyp.Fields.List {
		if field == nil || len(field.Names) == 0 {
			continue
		}

		if field.Doc == nil {
			continue
		}

		fieldMarkers := NewMarkerSet()

		for _, comment := range field.Doc.List {
			if marker := extractMarker(comment); marker.Value != "" {
				fieldMarkers.Insert(marker)
			}
		}

		fieldName := field.Names[0].Name
		results.insertStructFieldMarkers(sTyp, fieldName, fieldMarkers)
	}
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

// Marker represents a marker extracted from a comment on a declaration.
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
type MarkerSet map[string][]Marker

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
		ms[marker.Value] = append(ms[marker.Value], marker)
	}
}

// Has returns whether a marker with the value given is present in the
// MarkerSet.
func (ms MarkerSet) Has(value string) bool {
	_, ok := ms[value]
	return ok
}
