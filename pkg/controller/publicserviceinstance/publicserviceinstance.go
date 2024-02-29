package publicserviceinstance

import (
	"context"
	"fmt"
	"path"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type PublicServiceInstanceController struct {
	k8sClient     kubernetes.Interface
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewPublicServiceInstanceController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	pc := &PublicServiceInstanceController{}
	pc.workqueue = utils.NewWorkQueue(pc.Reconcile)
	pc.openappClient = openappHelper.OpenAPPClient
	pc.k8sClient = openappHelper.K8sClient

	openappHelper.PublicServiceInstanceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pc.workqueue.Add(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPublicServiceIns, ok := oldObj.(*v1alpha1.PublicServiceInstance)
			if !ok {
				return
			}
			newPublicServiceIns, ok := newObj.(*v1alpha1.PublicServiceInstance)
			if !ok {
				return
			}
			if reflect.DeepEqual(oldPublicServiceIns.Spec, newPublicServiceIns.Spec) {
				return
			}
			pc.workqueue.Add(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			pc.workqueue.Add(obj)
		},
	})

	return pc
}

func (pc *PublicServiceInstanceController) Start() {
	go pc.workqueue.Run()
}

func (pc *PublicServiceInstanceController) Reconcile(obj interface{}) error {
	publicServiceIns, ok := obj.(*v1alpha1.PublicServiceInstance)
	if !ok {
		klog.Errorf("Failed to convert obj to PublicServiceInstance")
		return nil
	}

	klog.Infof("Reconciling publicservice instance(%s/%s)...", publicServiceIns.Namespace, publicServiceIns.Name)
	if !publicServiceIns.DeletionTimestamp.IsZero() {
		return pc.deletePublicServiceInstanceResources(publicServiceIns)
	}

	publicServiceIns.Finalizers = []string{utils.PublicServiceControllerFinalizerKey}
	_, err := pc.openappClient.ServiceV1alpha1().PublicServiceInstances(publicServiceIns.Namespace).
		Update(context.Background(), publicServiceIns, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update publicservice instance finalizers: %v", err)
		return err
	}

	publicServiceTemp := publicServiceIns.Spec.PublicServiceTemplate
	if publicServiceTemp == "" {
		klog.Errorf("PublicServiceTemplate is not specified in PublicServiceInstance(%s/%s)", publicServiceIns.Namespace, publicServiceIns.Name)
		return nil
	}

	derivedResource := []commonv1alpha1.DerivedResource{}
	manifests := utils.FindTemplateResources(publicServiceTemp, utils.PublicServiceTemplateBasePath)
	for _, manifest := range manifests {
		if err := pc.handlePublicServiceInstanceDerivedResourceCreation(publicServiceIns, manifest, &derivedResource); err != nil {
			return err
		}
	}

	publicServiceIns.Status.DerivedResources = derivedResource
	_, err = pc.openappClient.ServiceV1alpha1().PublicServiceInstances(publicServiceIns.Namespace).
		UpdateStatus(context.Background(), publicServiceIns, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update publicservice instance status: %v", err)
		return err
	}

	return nil
}

func (pc *PublicServiceInstanceController) handlePublicServiceInstanceDerivedResourceCreation(pubclicServiceIns *v1alpha1.PublicServiceInstance,
	manifest string,
	derivedResource *[]commonv1alpha1.DerivedResource) error {
	values, err := utils.ConstructPublicServiceInstanceValues(pubclicServiceIns)
	if err != nil {
		return err
	}
	manifestContent, err := utils.ConstructTemplateWithValues(manifest, values)
	if err != nil {
		klog.Errorf("Failed to construct template with values: %v", err)
		return err
	}
	file := path.Base(manifest)
	switch file {
	case utils.TemplateManifestStatefulSetFile:
		err = pc.statefulsetHandler(manifestContent, derivedResource, pubclicServiceIns.Name)
		if err != nil {
			return err
		}
	case utils.TemplateManifestConfigMapFile:
		err = pc.configmapHandler(manifestContent, derivedResource, pubclicServiceIns.Name)
		if err != nil {
			return err
		}
	case utils.TemplateManifestServiceFile:
		err = pc.serviceHandler(manifestContent, derivedResource, pubclicServiceIns.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pc *PublicServiceInstanceController) deletePublicServiceInstanceResources(publicServiceIns *v1alpha1.PublicServiceInstance) error {
	klog.Infof("Deleting publicservice instance(%s/%s)...", publicServiceIns.Namespace, publicServiceIns.Name)

	appIns, err := pc.openappClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("Failed to list app instances: %v", err)
		return err
	}
	for _, ins := range appIns.Items {
		if ins.Spec.PublicServiceClass == publicServiceIns.Name {
			publicServiceIns.Status.Message = "APP instance is using this publicservice instance, cannot be deleted"
			pc.openappClient.ServiceV1alpha1().PublicServiceInstances(publicServiceIns.Namespace).
				UpdateStatus(context.Background(), publicServiceIns, metav1.UpdateOptions{})
			return fmt.Errorf("app instance is using this publicservice instance(%s/%s), cannot be deleted",
				publicServiceIns.Namespace, publicServiceIns.Name)
		}
	}

	for _, d := range publicServiceIns.Status.DerivedResources {
		switch d.Kind {
		case utils.InstanceDerivedResourceServiceKind:
			if err := utils.CleanInstanceDerivedServiceResource(pc.k8sClient,
				d.Name, publicServiceIns.Namespace); err != nil {
				return err
			}
		case utils.InstanceDerivedResourceConfigMapKind:
			if err := utils.CleanInstanceDerivedConfigMapResource(pc.k8sClient,
				d.Name, publicServiceIns.Namespace); err != nil {
				return err
			}
		case utils.InstanceDerivedResourceStatefulSetKind:
			if err := utils.CleanInstanceDerivedStatefulSetResource(pc.k8sClient,
				d.Name, publicServiceIns.Namespace); err != nil {
				return err
			}
		}
	}

	publicServiceIns.Finalizers = nil
	_, err = pc.openappClient.ServiceV1alpha1().PublicServiceInstances(publicServiceIns.Namespace).
		Update(context.Background(), publicServiceIns, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to remove publicservice instance finalizers: %v", err)
		return err
	}

	return nil
}

func (pc *PublicServiceInstanceController) serviceHandler(manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	instanceName string) error {
	labels := map[string]string{
		utils.PublicServiceInstanceLabelKey: instanceName,
	}
	return utils.CreateOrUpdateService(pc.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}

func (pc *PublicServiceInstanceController) configmapHandler(manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	instanceName string) error {
	labels := map[string]string{
		utils.PublicServiceInstanceLabelKey: instanceName,
	}
	return utils.CreateOrUpdateConfigmap(pc.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}

func (pc *PublicServiceInstanceController) statefulsetHandler(manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	instanceName string) error {
	labels := map[string]string{
		utils.PublicServiceInstanceLabelKey: instanceName,
	}
	return utils.CreateOrUpdateStatefulset(pc.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}
