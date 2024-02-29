package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Namespaced",categories={openapp-dev}
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels=openapp.dev/crd-install=true
// +kubebuilder:printcolumn:JSONPath=`.spec.publicServiceTemplate`,name=`PUBLIC-SERVICE-TEMPLATE`,type=string
// +kubebuilder:printcolumn:JSONPath=`.status.publicServiceReady`,name=`PUBLIC-SERVICE-READY`,type=string
// +kubebuilder:printcolumn:JSONPath=`.status.localServiceURL`,name=`LOCAL-URL`,type=string

type PublicServiceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec PublicServiceInstanceSpec `json:"spec"`

	// +optional
	Status PublicServiceInstanceStatus `json:"status"`
}

type PublicServiceInstanceSpec struct {
	PublicServiceTemplate string `json:"publicServiceTemplate"`
	Inputs                string `json:"inputs,omitempty"`
}

type PublicServiceInstanceStatus struct {
	// +optional
	PublicServiceReady bool `json:"publicServiceReady,omitempty"`
	// If there is service resource, the URL will exist.
	// +optional
	LocalServiceURL string `json:"localServiceURL,omitempty"`
	// +optional
	DerivedResources []commonv1alpha1.DerivedResource `json:"derivedResources,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PublicServiceInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PublicServiceInstance `json:"items"`
}
