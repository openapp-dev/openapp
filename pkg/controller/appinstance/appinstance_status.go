package appinstance

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type AppInstanceStatusController struct {
	k8sClient     kubernetes.Interface
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewAppInstanceStatusController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	ac := &AppInstanceStatusController{}
	ac.workqueue = utils.NewWorkQueue(ac.Reconcile)
	ac.openappClient = openappHelper.OpenAPPClient
	ac.k8sClient = openappHelper.K8sClient

	openappHelper.StatefulSetInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			sts, ok := obj.(*v1.StatefulSet)
			if !ok {
				return false
			}
			if sts.Namespace != utils.InstanceNamespace {
				return false
			}
			if sts.Labels == nil {
				return false
			}
			if _, ok := sts.Labels[utils.AppInstanceLabelKey]; !ok {
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
			DeleteFunc: func(obj interface{}) {
				ac.workqueue.Add(obj)
			},
		},
	})

	return ac
}

func (ac *AppInstanceStatusController) Start() {
	go ac.workqueue.Run()
}

func (ac *AppInstanceStatusController) Reconcile(obj interface{}) error {
	sts, ok := obj.(*v1.StatefulSet)
	if !ok {
		klog.Errorf("Failed to convert obj to StatefulSet")
		return nil
	}

	klog.Infof("Reconciling app instance status with statefulset(%s/%s)...", sts.Namespace, sts.Name)
	ready := false
	if sts.Status.ReadyReplicas > 0 {
		ready = true
	}

	instanceName := sts.Labels[utils.AppInstanceLabelKey]
	ins, err := ac.openappClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).
		Get(context.Background(), instanceName, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		return nil
	}

	insCopy := ins.DeepCopy()
	insCopy.Status.AppReady = ready
	_, err = ac.openappClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).
		UpdateStatus(context.Background(), insCopy, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update app instance status: %v", err)
		return err
	}

	return nil
}
