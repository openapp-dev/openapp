package publicserviceinstance

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgtypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type PublicServiceInstanceServiceController struct {
	k8sClient     kubernetes.Interface
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewPublicServiceInstanceServiceController(openappHandler *utils.OpenAPPHelper) types.ControllerInterface {
	pc := &PublicServiceInstanceServiceController{}
	pc.workqueue = utils.NewWorkQueue(pc.Reconcile)
	pc.k8sClient = openappHandler.K8sClient
	pc.openappClient = openappHandler.OpenAPPClient

	handlefunc := func(obj interface{}) {
		svc, ok := obj.(*corev1.Service)
		if !ok {
			return
		}
		pc.workqueue.Add(pkgtypes.NamespacedName{
			Namespace: utils.InstanceNamespace,
			Name:      svc.Name,
		})
	}

	_, _ = openappHandler.ServiceInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			svc, ok := obj.(*corev1.Service)
			if !ok {
				return false
			}
			if svc.Namespace != utils.InstanceNamespace {
				return false
			}
			if svc.Labels == nil {
				return false
			}
			if _, ok := svc.Labels[utils.PublicServiceInstanceLabelKey]; !ok {
				return false
			}
			return true
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				handlefunc(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				handlefunc(newObj)
			},
			DeleteFunc: func(obj interface{}) {
				handlefunc(obj)
			},
		},
	})

	return pc
}

func (pc *PublicServiceInstanceServiceController) Start() {
	go pc.workqueue.Run()
}

func (pc *PublicServiceInstanceServiceController) Reconcile(resourceKey pkgtypes.NamespacedName) error {
	klog.Infof("Reconciling publicservice with service(%s)...", resourceKey)
	svc, err := pc.k8sClient.CoreV1().Services(resourceKey.Namespace).
		Get(context.Background(), resourceKey.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return pc.updatePublicServiceInstanceServiceURL(svc, "")
		}
		klog.Errorf("Failed to get services:%v", err)
		return err
	}

	localURL, err := pc.getServiceURL(svc)
	if err != nil {
		klog.Errorf("Failed to get service url: %v", err)
		return err
	}

	return pc.updatePublicServiceInstanceServiceURL(svc, localURL)
}

func (sc *PublicServiceInstanceServiceController) getServiceURL(service *corev1.Service) (string, error) {
	localPort := service.Spec.Ports[0].NodePort
	url, err := utils.GetLocalServerIPAddress()
	if err != nil {
		klog.Errorf("Failed to get local server ip address: %v", err)
		return "", err
	}
	localURL := url + ":" + strconv.Itoa(int(localPort))
	return "http://" + localURL, nil
}

func (sc *PublicServiceInstanceServiceController) updatePublicServiceInstanceServiceURL(service *corev1.Service, localURL string) error {
	insName := service.Labels[utils.PublicServiceInstanceLabelKey]
	ins, err := sc.openappClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).Get(context.Background(), insName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("Failed to get publicservice instance: %v", err)
		return err
	}
	insCopy := ins.DeepCopy()
	insCopy.Status.LocalServiceURL = localURL
	_, err = sc.openappClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).UpdateStatus(context.Background(), insCopy, metav1.UpdateOptions{})

	return err
}
