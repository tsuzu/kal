package optionalorrequired

import (
	"go/ast"

	"github.com/JoelSpeed/kal/pkg/analysis/markers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"k8s.io/kube-openapi/pkg/util/sets"
)

const (
	// optionalMarker is the marker that indicates that a field is optional.
	optionalMarker = "optional"

	// requiredMarker is the marker that indicates that a field is required.
	requiredMarker = "required"

	// kubebuilderOptionalMarker is the marker that indicates that a field is optional in kubebuilder.
	kubebuilderOptionalMarker = "kubebuilder:validation:Optional"

	// kubebuilderRequiredMarker is the marker that indicates that a field is required in kubebuilder.
	kubebuilderRequiredMarker = "kubebuilder:validation:Required"
)

// Analyzer is the analyzer for the optionalorrequired package.
// It checks that all struct fields are marked either with the optional or required markers.
// It also checks that upstream markers are preferred over kubebuilder markers.
var Analyzer = &analysis.Analyzer{
	Name:     "optionalorrequired",
	Doc:      "Checks that all struct fields are marked either with the optional or required markers.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, markers.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	markers := pass.ResultOf[markers.Analyzer].(*markers.Markers)

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		sTyp, ok := n.(*ast.StructType)
		if !ok {
			return
		}

		structMarkers := markers.StructMarkers[sTyp]

		if sTyp.Fields == nil {
			return
		}

		for _, field := range sTyp.Fields.List {
			if field == nil || len(field.Names) == 0 {
				continue
			}

			fieldName := field.Names[0].Name
			fieldMarkers := structMarkers.FieldMarkers[fieldName]

			checkField(pass, field, fieldMarkers)
		}
	})

	return nil, nil
}

func checkField(pass *analysis.Pass, field *ast.Field, fieldMarkers []string) {
	if field == nil || len(field.Names) == 0 {
		return
	}

	fieldName := field.Names[0].Name

	markerSet := sets.NewString(fieldMarkers...)

	hasOptional := markerSet.Has(optionalMarker)
	hasRequired := markerSet.Has(requiredMarker)

	hasKubebuilderOptional := markerSet.Has(kubebuilderOptionalMarker)
	hasKubebuilderRequired := markerSet.Has(kubebuilderRequiredMarker)

	hasEitherOptional := hasOptional || hasKubebuilderOptional
	hasEitherRequired := hasRequired || hasKubebuilderRequired

	switch {
	case hasEitherOptional && hasEitherRequired:
		pass.Reportf(field.Pos(), "field %s must not be marked as both optional and required", fieldName)
	case hasKubebuilderOptional:
		pass.Reportf(field.Pos(), "field %s should use marker optional instead of kubebuilder:validation:Optional", fieldName)
	case hasKubebuilderRequired:
		pass.Reportf(field.Pos(), "field %s should use marker required instead of kubebuilder:validation:Required", fieldName)
	case hasOptional || hasRequired:
		// This is the correct state.
	default:
		pass.Reportf(field.Pos(), "field %s must be marked as optional or required", fieldName)
	}
}
