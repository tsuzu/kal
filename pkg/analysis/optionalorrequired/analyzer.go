package optionalorrequired

import (
	"fmt"
	"go/ast"

	"github.com/JoelSpeed/kal/pkg/analysis/markers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
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

func checkField(pass *analysis.Pass, field *ast.Field, fieldMarkers markers.MarkerSet) {
	if field == nil || len(field.Names) == 0 {
		return
	}

	fieldName := field.Names[0].Name

	hasOptional := fieldMarkers.Has(optionalMarker)
	hasRequired := fieldMarkers.Has(requiredMarker)

	hasKubebuilderOptional := fieldMarkers.Has(kubebuilderOptionalMarker)
	hasKubebuilderRequired := fieldMarkers.Has(kubebuilderRequiredMarker)

	hasEitherOptional := hasOptional || hasKubebuilderOptional
	hasEitherRequired := hasRequired || hasKubebuilderRequired

	hasBothOptional := hasOptional && hasKubebuilderOptional
	hasBothRequired := hasRequired && hasKubebuilderRequired

	switch {
	case hasEitherOptional && hasEitherRequired:
		pass.Reportf(field.Pos(), "field %s must not be marked as both optional and required", fieldName)
	case hasKubebuilderOptional:
		marker := fieldMarkers[kubebuilderOptionalMarker]
		if hasBothOptional {
			pass.Report(reportShouldRemoveKubebuilderMarker(field, marker, optionalMarker, kubebuilderOptionalMarker))
		} else {
			pass.Report(reportShouldReplaceKubebuilderMarker(field, marker, optionalMarker, kubebuilderOptionalMarker))
		}
	case hasKubebuilderRequired:
		marker := fieldMarkers[kubebuilderRequiredMarker]
		if hasBothRequired {
			pass.Report(reportShouldRemoveKubebuilderMarker(field, marker, requiredMarker, kubebuilderRequiredMarker))
		} else {
			pass.Report(reportShouldReplaceKubebuilderMarker(field, marker, requiredMarker, kubebuilderRequiredMarker))
		}
	case hasOptional || hasRequired:
		// This is the correct state.
	default:
		pass.Reportf(field.Pos(), "field %s must be marked as optional or required", fieldName)
	}
}

func reportShouldReplaceKubebuilderMarker(field *ast.Field, marker markers.Marker, desiredMarker, kubebuilderMaker string) analysis.Diagnostic {
	fieldName := field.Names[0].Name

	return analysis.Diagnostic{
		Pos:     field.Pos(),
		Message: fmt.Sprintf("field %s should use marker %s instead of %s", fieldName, desiredMarker, kubebuilderMaker),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("should replace `%s` with `%s`", kubebuilderMaker, desiredMarker),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     marker.Pos,
						End:     marker.End,
						NewText: []byte(fmt.Sprintf("// +%s", desiredMarker)),
					},
				},
			},
		},
	}
}

func reportShouldRemoveKubebuilderMarker(field *ast.Field, marker markers.Marker, desiredMarker, kubebuilderMaker string) analysis.Diagnostic {
	fieldName := field.Names[0].Name

	return analysis.Diagnostic{
		Pos:     field.Pos(),
		Message: fmt.Sprintf("field %s should use only the marker %s, %s is not required", fieldName, desiredMarker, kubebuilderMaker),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("should remove `// +%s`", kubebuilderMaker),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     marker.Pos,
						End:     marker.End + 1, // Add 1 to position to include the new line
						NewText: nil,
					},
				},
			},
		},
	}
}
