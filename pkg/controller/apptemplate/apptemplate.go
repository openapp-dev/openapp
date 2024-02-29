package apptemplate

import (
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	appv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type AppTemplateController struct {
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewAppTempalteController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	ac := &AppTemplateController{}
	ac.workqueue = utils.NewWorkQueue(ac.Reconcile)
	ac.openappClient = openappHelper.OpenAPPClient

	openappHelper.ConfigMapInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			cm := obj.(*v1.ConfigMap)
			if cm.Name != utils.SystemConfigMap || cm.Namespace != utils.SystemNamespace {
				return false
			}
			return true
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				ac.workqueue.Add(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				ac.workqueue.Add(newObj)
			},
		},
	})

	return ac
}

func (ac *AppTemplateController) Start() {
	go ac.workqueue.Run()
}

func (ac *AppTemplateController) Reconcile(_ interface{}) error {
	klog.Infof("Reconciling app template...")

	registries := utils.GetRegistryPaths()
	for _, registry := range registries {
		templates := utils.GetAppTemplatePath(registry)
		for _, templateFile := range templates {
			d, err := os.ReadFile(templateFile)
			if err != nil {
				klog.Errorf("Failed to read template file: %v", err)
				continue
			}
			appTemplate := &appv1alpha1.AppTemplate{}
			if err := yaml.Unmarshal(d, appTemplate); err != nil {
				klog.Errorf("Failed to unmarshal template: %v", err)
				continue
			}
			if err := createOrUpdateAppTemplate(ac.openappClient, appTemplate); err != nil {
				return err
			}
		}
	}

	return nil
}

func createOrUpdateAppTemplate(openappClient versioned.Interface, appTemplate *appv1alpha1.AppTemplate) error {
	templateExist, err := openappClient.AppV1alpha1().AppTemplates().Get(context.Background(), appTemplate.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("Failed to get app template: %v", err)
			return err
		}
		_, err = openappClient.AppV1alpha1().AppTemplates().Create(context.Background(), appTemplate, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("Failed to create app template: %v", err)
			return err
		}
		return nil
	}

	template := templateExist.DeepCopy()
	template.Spec = appTemplate.Spec
	_, err = openappClient.AppV1alpha1().AppTemplates().Update(context.Background(), template, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update app template: %v", err)
		return err
	}
	return err
}
