package publicserviceinstance

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

type PublicServiceInstanceStatusController struct {
	k8sClient     kubernetes.Interface
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewPublicServiceInstanceStatusController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	pc := &PublicServiceInstanceStatusController{}
	pc.workqueue = utils.NewWorkQueue(pc.Reconcile)
	pc.openappClient = openappHelper.OpenAPPClient
	pc.k8sClient = openappHelper.K8sClient

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
			if _, ok := sts.Labels[utils.PublicServiceInstanceLabelKey]; !ok {
				return false
			}
			return true
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pc.workqueue.Add(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pc.workqueue.Add(newObj)
			},
			DeleteFunc: func(obj interface{}) {
				pc.workqueue.Add(obj)
			},
		},
	})

	return pc
}

func (pc *PublicServiceInstanceStatusController) Start() {
	go pc.workqueue.Run()
}

func (pc *PublicServiceInstanceStatusController) Reconcile(obj interface{}) error {
	sts, ok := obj.(*v1.StatefulSet)
	if !ok {
		return nil
	}

	klog.Infof("Reconciling publicservice instance(%s/%s) status...", sts.Namespace, sts.Name)
	ready := false
	if sts.Status.ReadyReplicas > 0 {
		ready = true
	}

	instanceName := sts.Labels[utils.PublicServiceInstanceLabelKey]
	ins, err := pc.openappClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).
		Get(context.Background(), instanceName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}

		klog.Errorf("Failed to get publicservice instance: %v", err)
		return err
	}

	insCopy := ins.DeepCopy()
	insCopy.Status.PublicServiceReady = ready
	_, err = pc.openappClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).
		UpdateStatus(context.Background(), insCopy, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update publicservice instance status: %v", err)
		return err
	}

	return nil
}
