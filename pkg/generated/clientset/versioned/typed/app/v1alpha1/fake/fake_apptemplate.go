/*
Copyright 2024 The OpenAPP Authors.
SPDX-License-Identifier: BUSL-1.1
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAppTemplates implements AppTemplateInterface
type FakeAppTemplates struct {
	Fake *FakeAppV1alpha1
}

var apptemplatesResource = schema.GroupVersionResource{Group: "app.openapp.dev", Version: "v1alpha1", Resource: "apptemplates"}

var apptemplatesKind = schema.GroupVersionKind{Group: "app.openapp.dev", Version: "v1alpha1", Kind: "AppTemplate"}

// Get takes name of the appTemplate, and returns the corresponding appTemplate object, and an error if there is any.
func (c *FakeAppTemplates) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(apptemplatesResource, name), &v1alpha1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppTemplate), err
}

// List takes label and field selectors, and returns the list of AppTemplates that match those selectors.
func (c *FakeAppTemplates) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AppTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(apptemplatesResource, apptemplatesKind, opts), &v1alpha1.AppTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AppTemplateList{ListMeta: obj.(*v1alpha1.AppTemplateList).ListMeta}
	for _, item := range obj.(*v1alpha1.AppTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested appTemplates.
func (c *FakeAppTemplates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(apptemplatesResource, opts))
}

// Create takes the representation of a appTemplate and creates it.  Returns the server's representation of the appTemplate, and an error, if there is any.
func (c *FakeAppTemplates) Create(ctx context.Context, appTemplate *v1alpha1.AppTemplate, opts v1.CreateOptions) (result *v1alpha1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(apptemplatesResource, appTemplate), &v1alpha1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppTemplate), err
}

// Update takes the representation of a appTemplate and updates it. Returns the server's representation of the appTemplate, and an error, if there is any.
func (c *FakeAppTemplates) Update(ctx context.Context, appTemplate *v1alpha1.AppTemplate, opts v1.UpdateOptions) (result *v1alpha1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(apptemplatesResource, appTemplate), &v1alpha1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppTemplate), err
}

// Delete takes name of the appTemplate and deletes it. Returns an error if one occurs.
func (c *FakeAppTemplates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(apptemplatesResource, name, opts), &v1alpha1.AppTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAppTemplates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(apptemplatesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AppTemplateList{})
	return err
}

// Patch applies the patch and returns the patched appTemplate.
func (c *FakeAppTemplates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AppTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(apptemplatesResource, name, pt, data, subresources...), &v1alpha1.AppTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppTemplate), err
}
