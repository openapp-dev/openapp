package appinstance

import (
	"context"
	"path"
	"reflect"
	"strconv"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgtypes "k8s.io/apimachinery/pkg/types"
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

	_, _ = openappHelper.AppInstanceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ins, ok := obj.(*appv1alpha1.AppInstance)
			if !ok {
				return
			}
			ac.workqueue.Add(pkgtypes.NamespacedName{
				Namespace: utils.InstanceNamespace,
				Name:      ins.Name,
			})
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
			if newAppIns.DeletionTimestamp.IsZero() &&
				reflect.DeepEqual(oldAppIns.Spec, newAppIns.Spec) {
				return
			}
			ac.workqueue.Add(pkgtypes.NamespacedName{
				Namespace: utils.InstanceNamespace,
				Name:      newAppIns.Name,
			})
		},
		DeleteFunc: func(obj interface{}) {
			ins, ok := obj.(*appv1alpha1.AppInstance)
			if !ok {
				return
			}
			ac.workqueue.Add(pkgtypes.NamespacedName{
				Namespace: utils.InstanceNamespace,
				Name:      ins.Name,
			})
		},
	})

	return ac
}

func (ac *AppInstanceController) Start() {
	go ac.workqueue.Run()
}

func (ac *AppInstanceController) Reconcile(resourceKey pkgtypes.NamespacedName) error {
	klog.Infof("Reconciling app instance(%s)...", resourceKey)
	appIns, err := ac.openappClient.AppV1alpha1().AppInstances(resourceKey.Namespace).
		Get(context.Background(), resourceKey.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("Failed to get app instance: %v", err)
		return err
	}

	if !appIns.DeletionTimestamp.IsZero() {
		return ac.deleteAppInstanceResources(appIns)
	}

	appIns.Finalizers = []string{utils.AppInstanceControllerFinalizerKey}
	appIns, err = ac.openappClient.AppV1alpha1().AppInstances(appIns.Namespace).
		Update(context.Background(), appIns, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update app instance: %v", err)
		return err
	}

	appTemplate := appIns.Spec.AppTemplate
	if appTemplate == "" {
		klog.Errorf("AppTemplate is empty in AppInstance(%s/%s)", appIns.Namespace, appIns.Name)
		return nil
	}
	derivedResoruce := []commonv1alpha1.DerivedResource{}
	manifests := utils.FindAppTemplateResources(appTemplate)
	// The last element is statefulset, we should recreate it
	for _, manifest := range manifests {
		if err := ac.handleAppInstanceDerivedResourceCreation(appIns, manifest, &derivedResoruce); err != nil {
			return err
		}
	}
	appIns.Status.DerivedResources = derivedResoruce
	_, err = ac.openappClient.AppV1alpha1().AppInstances(appIns.Namespace).
		UpdateStatus(context.Background(), appIns, metav1.UpdateOptions{})
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
		err = ac.statefulsetHandler(manifestContent, derivedResoruce, appIns)
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

	appInstance.Finalizers = nil
	_, err := ac.openappClient.AppV1alpha1().AppInstances(appInstance.Namespace).
		Update(context.Background(), appInstance, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to remove app instance finalizers: %v", err)
		return err
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
	appIns *appv1alpha1.AppInstance) error {
	labels := map[string]string{
		utils.AppInstanceLabelKey:        appIns.Name,
		utils.InstanceGenerationLabelKey: strconv.Itoa(int(appIns.Generation)),
	}
	return utils.CreateOrUpdateStatefulset(ac.k8sClient,
		manifestContent,
		derivedResoruce,
		labels)
}
