package registry

import (
	"context"
	"path"
	"reflect"
	"strings"
	"time"

	gitv5 "github.com/go-git/go-git/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/controller/types"
	"github.com/openapp-dev/openapp/pkg/utils"
)

type RegistryController struct {
	k8sClient kubernetes.Interface
	cmLister  corev1.ConfigMapLister
	workqueue *utils.WorkQueue
}

func NewRegistryController(openappHelper *utils.OpenAPPHelper) types.ControllerInterface {
	rc := &RegistryController{
		cmLister: openappHelper.ConfigMapLister,
	}
	rc.workqueue = utils.NewWorkQueue(rc.Reconcile)
	rc.k8sClient = openappHelper.K8sClient
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
				rc.workqueue.Add(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCM, ok := oldObj.(*v1.ConfigMap)
				if !ok {
					return
				}
				newCM, ok := newObj.(*v1.ConfigMap)
				if !ok {
					return
				}
				if reflect.DeepEqual(oldCM.Data, newCM.Data) {
					return
				}
				rc.workqueue.Add(newObj)
			},
			DeleteFunc: func(obj interface{}) {
				rc.workqueue.Add(obj)
			},
		},
	})

	return rc
}

func (rc *RegistryController) Start() {
	go rc.workqueue.Run()

	ticker := time.NewTicker(time.Minute * 30)
	for range ticker.C {
		cm, err := rc.cmLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
		if err != nil {
			klog.Errorf("Failed to get config: %v", err)
			continue
		}
		rc.workqueue.Add(cm)
	}
}

func (rc *RegistryController) Reconcile(obj interface{}) error {
	klog.Infof("Reconciling config update...")

	cm, ok := obj.(*v1.ConfigMap)
	if !ok {
		klog.Errorf("Failed to convert object to configmap")
		return nil
	}

	registries := strings.Split(cm.Data[utils.RegistryKey], ",")
	if err := CloneOpenAPPRegistry(registries); err != nil {
		klog.Errorf("Failed to clone openapp registry: %v", err)
		return err
	}

	cmCopy := cm.DeepCopy()
	if cmCopy.Annotations == nil {
		cmCopy.Annotations = map[string]string{}
	}
	cmCopy.Annotations[utils.RegistryUpdateTimeAnnotationKey] = time.Now().Format(time.RFC3339)
	_, err := rc.k8sClient.CoreV1().ConfigMaps(utils.SystemNamespace).
		Update(context.Background(), cmCopy, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("Failed to update configmap: %v", err)
		return err
	}

	return nil
}

func CloneOpenAPPRegistry(registryList []string) error {
	for _, registry := range registryList {
		/// TBD: Can't clone specific branch, still don't know why
		repoURL, _, dir := getRepoURLAndBranch(registry)
		r, err := gitv5.PlainClone(path.Join(utils.RegistryCachePath, dir), false, &gitv5.CloneOptions{
			URL:             repoURL,
			InsecureSkipTLS: true,
		})
		if err != nil {
			if err != gitv5.ErrRepositoryAlreadyExists {
				klog.Errorf("Failed to clone registry %s: %v", registry, err)
				return err
			}
			r, err = gitv5.PlainOpen(path.Join(utils.RegistryCachePath, dir))
			if err != nil {
				klog.Errorf("Failed to open registry %s: %v", registry, err)
				return err
			}
		}
		w, err := r.Worktree()
		if err != nil {
			klog.Errorf("Failed to get worktree: %v", err)
			return err
		}
		// TBD: it couldn't pull the latest code sometimes
		err = w.Pull(&gitv5.PullOptions{RemoteName: "origin", InsecureSkipTLS: true, Force: true})
		if err != nil && err != gitv5.NoErrAlreadyUpToDate {
			klog.Errorf("Failed to pull registry %s: %v", registry, err)
			return err
		}
	}

	return nil
}

func getRepoURLAndBranch(url string) (string, string, string) {
	infos := strings.Split(url, "@")
	dirs := strings.Split(infos[0], "/")
	return infos[0], infos[1], dirs[len(dirs)-1]
}
