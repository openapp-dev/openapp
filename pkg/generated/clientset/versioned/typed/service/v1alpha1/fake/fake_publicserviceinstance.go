/*
Copyright 2024 The OpenAPP Authors.
SPDX-License-Identifier: BUSL-1.1
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePublicServiceInstances implements PublicServiceInstanceInterface
type FakePublicServiceInstances struct {
	Fake *FakeServiceV1alpha1
	ns   string
}

var publicserviceinstancesResource = schema.GroupVersionResource{Group: "service.openapp.dev", Version: "v1alpha1", Resource: "publicserviceinstances"}

var publicserviceinstancesKind = schema.GroupVersionKind{Group: "service.openapp.dev", Version: "v1alpha1", Kind: "PublicServiceInstance"}

// Get takes name of the publicServiceInstance, and returns the corresponding publicServiceInstance object, and an error if there is any.
func (c *FakePublicServiceInstances) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.PublicServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(publicserviceinstancesResource, c.ns, name), &v1alpha1.PublicServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PublicServiceInstance), err
}

// List takes label and field selectors, and returns the list of PublicServiceInstances that match those selectors.
func (c *FakePublicServiceInstances) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PublicServiceInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(publicserviceinstancesResource, publicserviceinstancesKind, c.ns, opts), &v1alpha1.PublicServiceInstanceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PublicServiceInstanceList{ListMeta: obj.(*v1alpha1.PublicServiceInstanceList).ListMeta}
	for _, item := range obj.(*v1alpha1.PublicServiceInstanceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested publicServiceInstances.
func (c *FakePublicServiceInstances) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(publicserviceinstancesResource, c.ns, opts))

}

// Create takes the representation of a publicServiceInstance and creates it.  Returns the server's representation of the publicServiceInstance, and an error, if there is any.
func (c *FakePublicServiceInstances) Create(ctx context.Context, publicServiceInstance *v1alpha1.PublicServiceInstance, opts v1.CreateOptions) (result *v1alpha1.PublicServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(publicserviceinstancesResource, c.ns, publicServiceInstance), &v1alpha1.PublicServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PublicServiceInstance), err
}

// Update takes the representation of a publicServiceInstance and updates it. Returns the server's representation of the publicServiceInstance, and an error, if there is any.
func (c *FakePublicServiceInstances) Update(ctx context.Context, publicServiceInstance *v1alpha1.PublicServiceInstance, opts v1.UpdateOptions) (result *v1alpha1.PublicServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(publicserviceinstancesResource, c.ns, publicServiceInstance), &v1alpha1.PublicServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PublicServiceInstance), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePublicServiceInstances) UpdateStatus(ctx context.Context, publicServiceInstance *v1alpha1.PublicServiceInstance, opts v1.UpdateOptions) (*v1alpha1.PublicServiceInstance, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(publicserviceinstancesResource, "status", c.ns, publicServiceInstance), &v1alpha1.PublicServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PublicServiceInstance), err
}

// Delete takes name of the publicServiceInstance and deletes it. Returns an error if one occurs.
func (c *FakePublicServiceInstances) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(publicserviceinstancesResource, c.ns, name, opts), &v1alpha1.PublicServiceInstance{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePublicServiceInstances) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(publicserviceinstancesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.PublicServiceInstanceList{})
	return err
}

// Patch applies the patch and returns the patched publicServiceInstance.
func (c *FakePublicServiceInstances) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.PublicServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(publicserviceinstancesResource, c.ns, name, pt, data, subresources...), &v1alpha1.PublicServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PublicServiceInstance), err
}
