package a

type OptionalOrRequiredTestStruct struct {
	NoMarkers string // want "field NoMarkers must be marked as optional or required"

	// noOptionalOrRequiredMarker is a field with no optional or required marker.
	// +enum
	// +kubebuilder:validation:Enum=Foo;Bar
	NoOptionalOrRequiredMarker string // want "field NoOptionalOrRequiredMarker must be marked as optional or required"

	// +optional
	// +required
	MarkedOpitonalAndRequired string // want "field MarkedOpitonalAndRequired must not be marked as both optional and required"

	// +optional
	// +kubebuilder:validation:Optional
	MarkedOptionalAndKubeBuilderOptional string // want "field MarkedOptionalAndKubeBuilderOptional should use only the marker optional, kubebuilder:validation:Optional is not required"

	// +required
	// +kubebuilder:validation:Required
	MarkedRequiredAndKubeBuilderRequired string // want "field MarkedRequiredAndKubeBuilderRequired should use only the marker required, kubebuilder:validation:Required is not required"

	// +kubebuilder:validation:Optional
	KubebuilderOptionalMarker string // want "field KubebuilderOptionalMarker should use marker optional instead of kubebuilder:validation:Optional"

	// +kubebuilder:validation:Required
	KubebuilderRequiredMarker string // want "field KubebuilderRequiredMarker should use marker required instead of kubebuilder:validation:Required"

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Optional
	MarkedWithKubeBuilderOptionalTwice string // want "field MarkedWithKubeBuilderOptionalTwice should use marker optional instead of kubebuilder:validation:Optional"

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Optional
	// +optional
	MarkedWithKubeBuilderOptionalTwiceAndOptional string // want "field MarkedWithKubeBuilderOptionalTwiceAndOptional should use only the marker optional, kubebuilder:validation:Optional is not required"

	// +optional
	OptionalMarker string

	// +required
	RequiredMarker string
}
