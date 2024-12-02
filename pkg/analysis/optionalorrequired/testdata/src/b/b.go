package a

type OptionalOrRequiredTestStruct struct {
	NoMarkers string // want "field NoMarkers must be marked as kubebuilder:validation:Optional or kubebuilder:validation:Required"

	// noOptionalOrRequiredMarker is a field with no optional or required marker.
	// +enum
	// +kubebuilder:validation:Enum=Foo;Bar
	NoOptionalOrRequiredMarker string // want "field NoOptionalOrRequiredMarker must be marked as kubebuilder:validation:Optional or kubebuilder:validation:Required"

	// +optional
	// +required
	MarkedOpitonalAndRequired string // want "field MarkedOpitonalAndRequired must not be marked as both optional and required"

	// +optional
	// +kubebuilder:validation:Optional
	MarkedOptionalAndKubeBuilderOptional string // want "field MarkedOptionalAndKubeBuilderOptional should use only the marker kubebuilder:validation:Optional, optional is not required"

	// +required
	// +kubebuilder:validation:Required
	MarkedRequiredAndKubeBuilderRequired string // want "field MarkedRequiredAndKubeBuilderRequired should use only the marker kubebuilder:validation:Required, required is not required"

	// +kubebuilder:validation:Optional
	KubebuilderOptionalMarker string

	// +kubebuilder:validation:Required
	KubebuilderRequiredMarker string

	// +optional
	OptionalMarker string // want "field OptionalMarker should use marker kubebuilder:validation:Optional instead of optional"

	// +required
	RequiredMarker string // want "field RequiredMarker should use marker kubebuilder:validation:Required instead of required"

}
