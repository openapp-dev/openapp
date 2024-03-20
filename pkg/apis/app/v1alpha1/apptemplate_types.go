package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster",categories={openapp-dev}
// +kubebuilder:metadata:labels=openapp.dev/crd-install=true
// +kubebuilder:printcolumn:JSONPath=`.spec.url`,name=`APP-TEMPLATE-URL`,type=string
// +kubebuilder:printcolumn:JSONPath=`.spec.exposeType`,name=`EXPOSE-TYPE`,type=string

type AppTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AppTemplateSpec `json:"spec"`
}

type AppTemplateSpec struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Author      string                    `json:"author"`
	Icon        string                    `json:"icon"`
	URL         string                    `json:"url"`
	Inputs      string                    `json:"inputs"`
	ExposeType  commonv1alpha1.ExposeType `json:"exposeType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AppTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppTemplate `json:"items"`
}
