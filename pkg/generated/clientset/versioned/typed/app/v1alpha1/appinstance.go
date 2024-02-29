/*
Copyright 2024 The OpenAPP Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
	scheme "github.com/openapp-dev/openapp/pkg/generated/clientset/versioned/scheme"
)

// AppInstancesGetter has a method to return a AppInstanceInterface.
// A group's client should implement this interface.
type AppInstancesGetter interface {
	AppInstances(namespace string) AppInstanceInterface
}

// AppInstanceInterface has methods to work with AppInstance resources.
type AppInstanceInterface interface {
	Create(ctx context.Context, appInstance *v1alpha1.AppInstance, opts v1.CreateOptions) (*v1alpha1.AppInstance, error)
	Update(ctx context.Context, appInstance *v1alpha1.AppInstance, opts v1.UpdateOptions) (*v1alpha1.AppInstance, error)
	UpdateStatus(ctx context.Context, appInstance *v1alpha1.AppInstance, opts v1.UpdateOptions) (*v1alpha1.AppInstance, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.AppInstance, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AppInstanceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AppInstance, err error)
	AppInstanceExpansion
}

// appInstances implements AppInstanceInterface
type appInstances struct {
	client rest.Interface
	ns     string
}

// newAppInstances returns a AppInstances
func newAppInstances(c *AppV1alpha1Client, namespace string) *appInstances {
	return &appInstances{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the appInstance, and returns the corresponding appInstance object, and an error if there is any.
func (c *appInstances) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AppInstance, err error) {
	result = &v1alpha1.AppInstance{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("appinstances").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AppInstances that match those selectors.
func (c *appInstances) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AppInstanceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AppInstanceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("appinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested appInstances.
func (c *appInstances) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("appinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a appInstance and creates it.  Returns the server's representation of the appInstance, and an error, if there is any.
func (c *appInstances) Create(ctx context.Context, appInstance *v1alpha1.AppInstance, opts v1.CreateOptions) (result *v1alpha1.AppInstance, err error) {
	result = &v1alpha1.AppInstance{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("appinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(appInstance).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a appInstance and updates it. Returns the server's representation of the appInstance, and an error, if there is any.
func (c *appInstances) Update(ctx context.Context, appInstance *v1alpha1.AppInstance, opts v1.UpdateOptions) (result *v1alpha1.AppInstance, err error) {
	result = &v1alpha1.AppInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("appinstances").
		Name(appInstance.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(appInstance).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *appInstances) UpdateStatus(ctx context.Context, appInstance *v1alpha1.AppInstance, opts v1.UpdateOptions) (result *v1alpha1.AppInstance, err error) {
	result = &v1alpha1.AppInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("appinstances").
		Name(appInstance.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(appInstance).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the appInstance and deletes it. Returns an error if one occurs.
func (c *appInstances) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("appinstances").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *appInstances) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("appinstances").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched appInstance.
func (c *appInstances) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AppInstance, err error) {
	result = &v1alpha1.AppInstance{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("appinstances").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
