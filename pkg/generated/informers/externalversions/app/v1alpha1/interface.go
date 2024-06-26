/*
Copyright 2024 The OpenAPP Authors.
SPDX-License-Identifier: BUSL-1.1
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/openapp-dev/openapp/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// AppInstances returns a AppInstanceInformer.
	AppInstances() AppInstanceInformer
	// AppTemplates returns a AppTemplateInformer.
	AppTemplates() AppTemplateInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// AppInstances returns a AppInstanceInformer.
func (v *version) AppInstances() AppInstanceInformer {
	return &appInstanceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// AppTemplates returns a AppTemplateInformer.
func (v *version) AppTemplates() AppTemplateInformer {
	return &appTemplateInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
