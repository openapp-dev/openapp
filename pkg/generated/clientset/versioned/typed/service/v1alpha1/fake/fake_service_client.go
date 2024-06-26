/*
Copyright 2024 The OpenAPP Authors.
SPDX-License-Identifier: BUSL-1.1
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/openapp-dev/openapp/pkg/generated/clientset/versioned/typed/service/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeServiceV1alpha1 struct {
	*testing.Fake
}

func (c *FakeServiceV1alpha1) PublicServiceInstances(namespace string) v1alpha1.PublicServiceInstanceInterface {
	return &FakePublicServiceInstances{c, namespace}
}

func (c *FakeServiceV1alpha1) PublicServiceTemplates() v1alpha1.PublicServiceTemplateInterface {
	return &FakePublicServiceTemplates{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeServiceV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
