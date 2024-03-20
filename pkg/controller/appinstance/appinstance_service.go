package appinstance

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

	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type AppInstanceServiceController struct {
	k8sClient     kubernetes.Interface
	openappClient versioned.Interface
	workqueue     *utils.WorkQueue
}

func NewAppInstanceServiceController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	sc := &AppInstanceServiceController{}
	sc.workqueue = utils.NewWorkQueue(sc.Reconcile)
	sc.k8sClient = openappHelper.K8sClient
	sc.openappClient = openappHelper.OpenAPPClient

	handlefunc := func(obj interface{}) {
		svc, ok := obj.(*corev1.Service)
		if !ok {
			return
		}
		sc.workqueue.Add(pkgtypes.NamespacedName{
			Namespace: utils.InstanceNamespace,
			Name:      svc.Name,
		})
	}

	_, _ = openappHelper.ServiceInformer.AddEventHandler(cache.FilteringResourceEventHandler{
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
			if _, ok := svc.Labels[utils.AppInstanceLabelKey]; !ok {
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

	return sc
}

func (ac *AppInstanceServiceController) Start() {
	go ac.workqueue.Run()
}

func (sc *AppInstanceServiceController) Reconcile(resourceKey pkgtypes.NamespacedName) error {
	klog.Infof("Reconciling app instance service status with service(%s)...", resourceKey)
	svc, err := sc.k8sClient.CoreV1().Services(resourceKey.Namespace).
		Get(context.Background(), resourceKey.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return sc.updateAppInstanceServiceURL(svc, "", "")
		}
		klog.Errorf("Failed to get service: %v", err)
		return err
	}
	localURL, publicURL, err := sc.getServiceURL(svc)
	if err != nil {
		klog.Errorf("Failed to get service url: %v", err)
		return err
	}
	return sc.updateAppInstanceServiceURL(svc, localURL, publicURL)
}

func (sc *AppInstanceServiceController) updateAppInstanceServiceURL(service *corev1.Service, localURL, publicURL string) error {
	appInsName := service.Labels[utils.AppInstanceLabelKey]
	appIns, err := sc.openappClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).Get(context.Background(), appInsName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("Failed to get app instance: %v", err)
		return err
	}
	appInsCopy := appIns.DeepCopy()
	appInsCopy.Status.ExternalServiceURL = publicURL
	appInsCopy.Status.LocalServiceURL = localURL
	_, err = sc.openappClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).UpdateStatus(context.Background(), appInsCopy, metav1.UpdateOptions{})

	return err
}

func (sc *AppInstanceServiceController) getServiceURL(service *corev1.Service) (string, string, error) {
	appInsName := service.Labels[utils.AppInstanceLabelKey]
	appIns, err := sc.openappClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).Get(context.Background(), appInsName, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("Failed to get app instance: %v", err)
		return "", "", err
	}

	publicURL := ""
	if len(service.Status.LoadBalancer.Ingress) != 0 {
		port := int32(80)
		if len(service.Status.LoadBalancer.Ingress[0].Ports) != 0 {
			port = service.Status.LoadBalancer.Ingress[0].Ports[0].Port
		}
		url := service.Status.LoadBalancer.Ingress[0].Hostname
		if url == "" {
			url = service.Status.LoadBalancer.Ingress[0].IP
		}
		publicURL = url + ":" + strconv.Itoa(int(port))
	}

	localPort := service.Spec.Ports[0].NodePort
	url, err := utils.GetLocalServerIPAddress()
	if err != nil {
		klog.Errorf("Failed to get local server ip address: %v", err)
		return "", "", err
	}
	localURL := url + ":" + strconv.Itoa(int(localPort))

	appTemp, err := sc.openappClient.AppV1alpha1().AppTemplates().Get(context.Background(), appIns.Spec.AppTemplate, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("Failed to get app template: %v", err)
		return "", "", err
	}
	if appTemp.Spec.ExposeType == commonv1alpha1.ExposeLayer4 {
		return localURL, publicURL, nil
	}
	if publicURL == "" {
		return "http://" + localURL, "", nil
	}
	return "http://" + localURL, "http://" + publicURL, nil
}
