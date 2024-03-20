package v1alpha1

type DerivedResource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

type ExposeType string

const (
	ExposeLayer4 ExposeType = "Layer4"
	ExposeLayer7 ExposeType = "Layer7"
)
