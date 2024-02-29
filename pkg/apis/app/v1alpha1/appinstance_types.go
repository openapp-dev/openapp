package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={openapp-dev}
// +kubebuilder:metadata:labels=openapp.dev/crd-install=true
// +kubebuilder:printcolumn:JSONPath=`.spec.appTemplate`,name=`APP-TEMPLATE`,type=string
// +kubebuilder:printcolumn:JSONPath=`.status.appReady`,name=`APP-READY`,type=string
// +kubebuilder:printcolumn:JSONPath=`.spec.publicServiceClass`,name=`PUBLIC-SERVICE`,type=string
// +kubebuilder:printcolumn:JSONPath=`.status.externalServiceURL`,name=`PUBLIC-URL`,type=string
// +kubebuilder:printcolumn:JSONPath=`.status.localServiceURL`,name=`LOCAL-URL`,type=string

type AppInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec AppInstanceSpec `json:"spec"`

	// +optional
	Status AppInstanceStatus `json:"status"`
}

type AppInstanceSpec struct {
	PublicServiceClass string `json:"publicServiceClass,omitempty"`
	AppTemplate        string `json:"appTemplate"`
	Inputs             string `json:"inputs,omitempty"`
}

type AppInstanceStatus struct {
	// +optional
	AppReady bool `json:"appReady,omitempty"`
	// +optional
	ExternalServiceURL string `json:"externalServiceURL,omitempty"`
	// +optional
	LocalServiceURL string `json:"localServiceURL,omitempty"`
	// +optional
	DerivedResources []commonv1alpha1.DerivedResource `json:"derivedResources,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AppInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppInstance `json:"items"`
}
