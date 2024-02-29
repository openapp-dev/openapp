package appinstance

import (
	"context"
	"path"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	appv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type AppInstanceController struct {
	k8sClient     kubernetes.Interface
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewAppInstanceController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	ac := &AppInstanceController{}
	ac.workqueue = utils.NewWorkQueue(ac.Reconcile)
	ac.openappClient = openappHelper.OpenAPPClient
	ac.k8sClient = openappHelper.K8sClient

	openappHelper.AppInstanceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ac.workqueue.Add(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldAppIns, ok := oldObj.(*appv1alpha1.AppInstance)
			if !ok {
				return
			}
			newAppIns, ok := newObj.(*appv1alpha1.AppInstance)
			if !ok {
				return
			}
			if reflect.DeepEqual(oldAppIns.Spec, newAppIns.Spec) {
				return
			}
			ac.workqueue.Add(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			ac.workqueue.Add(obj)
		},
	})

	return ac
}

func (ac *AppInstanceController) Start() {
	go ac.workqueue.Run()
}

func (ac *AppInstanceController) Reconcile(obj interface{}) error {
	appInstance, ok := obj.(*appv1alpha1.AppInstance)
	if !ok {
		klog.Errorf("Failed to convert obj to AppInstance")
		return nil
	}

	klog.Infof("Reconciling app instance(%s/%s)...", appInstance.Namespace, appInstance.Name)
	appInsExists, err := ac.openappClient.AppV1alpha1().AppInstances(appInstance.Namespace).
		Get(context.Background(), appInstance.Name, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		return ac.deleteAppInstanceResources(appInstance)
	}

	appTemplate := appInstance.Spec.AppTemplate
	if appTemplate == "" {
		klog.Errorf("AppTemplate is empty in AppInstance(%s/%s)", appInstance.Namespace, appInstance.Name)
		return nil
	}

	derivedResoruce := []commonv1alpha1.DerivedResource{}
	manifests := utils.FindAppTemplateResources(appTemplate)
	// The last element is statefulset, we should recreate it
	for _, manifest := range manifests {
		if err := ac.handleAppInstanceDerivedResourceCreation(appInstance, manifest, &derivedResoruce); err != nil {
			return err
		}
	}
	appInsExists.Status.DerivedResources = derivedResoruce
	_, err = ac.openappClient.AppV1alpha1().AppInstances(appInstance.Namespace).
		UpdateStatus(context.Background(), appInsExists, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update app instance status: %v", err)
		return err
	}

	return nil
}

func (ac *AppInstanceController) handleAppInstanceDerivedResourceCreation(appIns *appv1alpha1.AppInstance,
	manifest string,
	derivedResoruce *[]commonv1alpha1.DerivedResource) error {
	values, err := utils.ConstructAppInstanceValues(appIns)
	if err != nil {
		return err
	}
	manifestContent, err := utils.ConstructTemplateWithValues(manifest, values)
	if err != nil {
		klog.Errorf("Failed to construct manifest with values: %v", err)
		return err
	}
	file := path.Base(manifest)
	switch file {
	case utils.TemplateManifestStatefulSetFile:
		err = ac.statefulsetHandler(manifestContent, derivedResoruce, appIns.Name)
		if err != nil {
			klog.Errorf("Failed to handle statefulset: %v", err)
			return err
		}
	// TBD: except statefulset, there is a general way to update the resource
	case utils.TemplateManifestServiceFile:
		err = ac.serviceHandler(manifestContent, derivedResoruce, appIns)
		if err != nil {
			klog.Errorf("Failed to handle service: %v", err)
			return err
		}
	case utils.TemplateManifestConfigMapFile:
		err = ac.configmapHandler(manifestContent, derivedResoruce, appIns.Name)
		if err != nil {
			klog.Errorf("Failed to handle configmap: %v", err)
			return err
		}
	}
	return nil
}

func (ac *AppInstanceController) deleteAppInstanceResources(appInstance *appv1alpha1.AppInstance) error {
	klog.Infof("Deleting app instance(%s/%s) resources...", appInstance.Namespace, appInstance.Name)
	for _, d := range appInstance.Status.DerivedResources {
		switch d.Kind {
		case utils.InstanceDerivedResourceServiceKind:
			if err := utils.CleanInstanceDerivedServiceResource(ac.k8sClient,
				d.Name, appInstance.Namespace); err != nil {
				return err
			}
		case utils.InstanceDerivedResourceConfigMapKind:
			if err := utils.CleanInstanceDerivedConfigMapResource(ac.k8sClient,
				d.Name, appInstance.Namespace); err != nil {
				return err
			}
		case utils.InstanceDerivedResourceStatefulSetKind:
			if err := utils.CleanInstanceDerivedStatefulSetResource(ac.k8sClient,
				d.Name, appInstance.Namespace); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ac *AppInstanceController) serviceHandler(manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	appIns *appv1alpha1.AppInstance) error {
	labels := map[string]string{
		utils.ServiceExposeClassLabelKey: appIns.Spec.PublicServiceClass,
		utils.AppInstanceLabelKey:        appIns.Name,
	}
	return utils.CreateOrUpdateService(ac.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}

func (ac *AppInstanceController) configmapHandler(manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	instanceName string) error {
	labels := map[string]string{
		utils.AppInstanceLabelKey: instanceName,
	}
	return utils.CreateOrUpdateConfigmap(ac.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}

func (ac *AppInstanceController) statefulsetHandler(manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	instanceName string) error {
	labels := map[string]string{
		utils.AppInstanceLabelKey: instanceName,
	}
	return utils.CreateOrUpdateStatefulset(ac.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}
