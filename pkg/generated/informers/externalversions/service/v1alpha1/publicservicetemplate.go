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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	servicev1alpha1 "github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
	versioned "github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/openapp-dev/openapp/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/openapp-dev/openapp/pkg/generated/listers/service/v1alpha1"
)

// PublicServiceTemplateInformer provides access to a shared informer and lister for
// PublicServiceTemplates.
type PublicServiceTemplateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PublicServiceTemplateLister
}

type publicServiceTemplateInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewPublicServiceTemplateInformer constructs a new informer for PublicServiceTemplate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPublicServiceTemplateInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPublicServiceTemplateInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredPublicServiceTemplateInformer constructs a new informer for PublicServiceTemplate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPublicServiceTemplateInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ServiceV1alpha1().PublicServiceTemplates().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ServiceV1alpha1().PublicServiceTemplates().Watch(context.TODO(), options)
			},
		},
		&servicev1alpha1.PublicServiceTemplate{},
		resyncPeriod,
		indexers,
	)
}

func (f *publicServiceTemplateInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPublicServiceTemplateInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *publicServiceTemplateInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&servicev1alpha1.PublicServiceTemplate{}, f.defaultInformer)
}

func (f *publicServiceTemplateInformer) Lister() v1alpha1.PublicServiceTemplateLister {
	return v1alpha1.NewPublicServiceTemplateLister(f.Informer().GetIndexer())
}
