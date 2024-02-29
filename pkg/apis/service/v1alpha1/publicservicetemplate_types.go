package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster",categories={openapp-dev}
// +kubebuilder:metadata:labels=openapp.dev/crd-install=true

type PublicServiceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec PublicServiceTemplateSpec `json:"spec"`
}

type PublicServiceTemplateSpec struct {
	// +required
	Title string `json:"title"`
	// +required
	Description string `json:"description"`
	Author      string `json:"author"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	Inputs      string `json:"inputs"`
	// +required
	ExposeTypes []ExposeType `json:"exposeTypes"`
}

type ExposeType string

const (
	ExposeLayer4 ExposeType = "Layer4"
	ExposeLayer7 ExposeType = "Layer7"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PublicServiceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PublicServiceTemplate `json:"items"`
}
