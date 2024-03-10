package publicservicetemplate

import (
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgtypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	servicev1alpha1 "github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type PublicServiceTemplateController struct {
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewPublicServiceTemplateController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	pc := &PublicServiceTemplateController{}
	pc.workqueue = utils.NewWorkQueue(pc.Reconcile)
	pc.openappClient = openappHelper.OpenAPPClient

	_, _ = openappHelper.ConfigMapInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			cm := obj.(*v1.ConfigMap)
			if cm.Name != utils.SystemConfigMap || cm.Namespace != utils.SystemNamespace {
				return false
			}
			return true
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(_ interface{}) {
				pc.workqueue.Add(pkgtypes.NamespacedName{})
			},
			UpdateFunc: func(_, _ interface{}) {
				pc.workqueue.Add(pkgtypes.NamespacedName{})
			},
		},
	})

	return pc
}

func (ac *PublicServiceTemplateController) Start() {
	go ac.workqueue.Run()
}

func (ac *PublicServiceTemplateController) Reconcile(_ pkgtypes.NamespacedName) error {
	klog.Infof("Reconciling publicservice template...")
	registries := utils.GetRegistryPaths()
	for _, registry := range registries {
		templates := utils.GetPublicServiceTemplatePath(registry)
		for _, templateFile := range templates {
			d, err := os.ReadFile(templateFile)
			if err != nil {
				klog.Errorf("Failed to read template file: %v", err)
				continue
			}
			serviceTemplate := &servicev1alpha1.PublicServiceTemplate{}
			if err := yaml.Unmarshal(d, serviceTemplate); err != nil {
				klog.Errorf("Failed to unmarshal template: %v", err)
				continue
			}
			if err := createOrUpdateServiceTemplate(ac.openappClient, serviceTemplate); err != nil {
				return err
			}
		}
	}
	return nil
}

func createOrUpdateServiceTemplate(openappClient versioned.Interface, serviceTemplate *servicev1alpha1.PublicServiceTemplate) error {
	templateExist, err := openappClient.ServiceV1alpha1().PublicServiceTemplates().Get(context.Background(), serviceTemplate.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("Failed to get public service template: %v", err)
			return err
		}
		_, err = openappClient.ServiceV1alpha1().PublicServiceTemplates().Create(context.Background(), serviceTemplate, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("Failed to create public service template: %v", err)
			return err
		}
		return nil
	}

	template := templateExist.DeepCopy()
	template.Spec = serviceTemplate.Spec
	_, err = openappClient.ServiceV1alpha1().PublicServiceTemplates().Update(context.Background(), template, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update public service template: %v", err)
		return err
	}
	return err
}
