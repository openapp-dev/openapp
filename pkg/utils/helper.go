package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	appv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
	commonv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/common/v1alpha1"
	servicev1alpha1 "github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	openappinformer "github.com/openapp-dev/openapp/pkg/generated/informers/externalversions"
	listerappv1alpha1 "github.com/openapp-dev/openapp/pkg/generated/listers/app/v1alpha1"
	listerservicev1alpha1 "github.com/openapp-dev/openapp/pkg/generated/listers/service/v1alpha1"
)

type OpenAPPHelper struct {
	K8sClient                     kubernetes.Interface
	OpenAPPClient                 versioned.Interface
	ConfigMapInformer             cache.SharedIndexInformer
	ServiceInformer               cache.SharedIndexInformer
	AppInstanceInformer           cache.SharedIndexInformer
	PublicServiceInstanceInformer cache.SharedIndexInformer
	StatefulSetInformer           cache.SharedIndexInformer
	ConfigMapLister               corev1.ConfigMapLister
	AppInstanceLister             listerappv1alpha1.AppInstanceLister
	AppTemplateLister             listerappv1alpha1.AppTemplateLister
	PublicServiceInstanceLister   listerservicev1alpha1.PublicServiceInstanceLister
	PublicServiceTemplateLister   listerservicev1alpha1.PublicServiceTemplateLister
}

func NewOpenAPPHelper(ctx context.Context,
	k8sClient kubernetes.Interface,
	openappClient versioned.Interface) *OpenAPPHelper {
	k8sFactory := informers.NewSharedInformerFactory(k8sClient, 0)
	openappFactory := openappinformer.NewSharedInformerFactory(openappClient, 0)

	configMapInformer := k8sFactory.Core().V1().ConfigMaps().Informer()
	serviceInformer := k8sFactory.Core().V1().Services().Informer()
	statefulSetInformer := k8sFactory.Apps().V1().StatefulSets().Informer()
	appInstanceInformer := openappFactory.App().V1alpha1().AppInstances().Informer()
	serviceInstanceInformer := openappFactory.Service().V1alpha1().PublicServiceInstances().Informer()

	helper := OpenAPPHelper{
		K8sClient:                     k8sClient,
		OpenAPPClient:                 openappClient,
		ConfigMapInformer:             configMapInformer,
		ServiceInformer:               serviceInformer,
		AppInstanceInformer:           appInstanceInformer,
		PublicServiceInstanceInformer: serviceInstanceInformer,
		StatefulSetInformer:           statefulSetInformer,
		ConfigMapLister:               k8sFactory.Core().V1().ConfigMaps().Lister(),
		AppInstanceLister:             openappFactory.App().V1alpha1().AppInstances().Lister(),
		AppTemplateLister:             openappFactory.App().V1alpha1().AppTemplates().Lister(),
		PublicServiceInstanceLister:   openappFactory.Service().V1alpha1().PublicServiceInstances().Lister(),
		PublicServiceTemplateLister:   openappFactory.Service().V1alpha1().PublicServiceTemplates().Lister(),
	}

	k8sFactory.Start(ctx.Done())
	openappFactory.Start(ctx.Done())

	return &helper
}

func GetRegistryPaths() []string {
	ret := []string{}
	dirs, err := os.ReadDir(RegistryCachePath)
	if err != nil {
		klog.Errorf("Failed to read registry cache path: %v", err)
		return ret
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			ret = append(ret, filepath.Join(RegistryCachePath, dir.Name()))
		}
	}
	return ret
}

func GetAppTemplatePath(registryPath string) []string {
	ret := []string{}
	basicPath := path.Join(registryPath, AppTemplatePath)
	dirs, err := os.ReadDir(basicPath)
	if err != nil {
		klog.Errorf("Failed to read registry app templa path: %v", err)
		return ret
	}
	for _, dir := range dirs {
		ret = append(ret, path.Join(basicPath, dir.Name(), TemplateFileName))
	}
	return ret
}

func FindAppTemplateResources(appTemplate string) []string {
	return FindTemplateResources(appTemplate, AppTemplateBasePath)
}

func FindTemplateResources(tempName, tempBasePath string) []string {
	registries := GetRegistryPaths()
	for _, regiregistry := range registries {
		templates := getTemplates(regiregistry, tempBasePath)
		for _, t := range templates {
			if t != tempName {
				continue
			}
			return getTemplateManifests(regiregistry, tempName, tempBasePath)
		}
	}
	return nil
}

func getTemplates(registryPath, tempBasePath string) []string {
	ret := []string{}
	basicPath := path.Join(registryPath, tempBasePath)
	dirs, err := os.ReadDir(basicPath)
	if err != nil {
		klog.Errorf("Failed to read registry templa path: %v", err)
		return ret
	}
	for _, dir := range dirs {
		ret = append(ret, dir.Name())
	}
	return ret
}

func getTemplateManifests(registry, tempName, tempBasePath string) []string {
	basicPath := path.Join(registry, tempBasePath, tempName, TemplateManifestsDirName)
	files, err := os.ReadDir(basicPath)
	if err != nil {
		klog.Errorf("Failed to read registry template manifests path: %v", err)
		return nil
	}
	ret := []string{}
	for _, f := range files {
		ret = append(ret, path.Join(basicPath, f.Name()))
	}
	sort.Slice(ret, func(i, j int) bool {
		return strings.HasSuffix(ret[i], "statefulset.yaml")
	})
	return ret
}

func GetPublicServiceTemplatePath(registryPath string) []string {
	ret := []string{}
	basicPath := path.Join(registryPath, PublicServiceTemplatePath)
	dirs, err := os.ReadDir(basicPath)
	if err != nil {
		klog.Errorf("Failed to read registry app templa path: %v", err)
		return ret
	}
	for _, dir := range dirs {
		ret = append(ret, path.Join(basicPath, dir.Name(), TemplateFileName))
	}
	return ret
}

func ConstructTemplateWithValues(manifestFile, instanceValues string) ([]byte, error) {
	var values map[string]interface{}
	err := json.Unmarshal([]byte(instanceValues), &values)
	if err != nil {
		klog.Errorf("Failed to unmarshal values: %v", err)
		return nil, err
	}
	// TDB: The parameter of New should be the file name, still don't know why
	tmpl, err := template.New(path.Base(manifestFile)).ParseFiles(manifestFile)
	if err != nil {
		klog.Errorf("Failed to parse template: %v", err)
		return nil, err
	}

	var ret bytes.Buffer
	err = tmpl.Execute(&ret, values)
	if err != nil {
		klog.Errorf("Failed to execute template: %v", err)
		return nil, err
	}
	return ret.Bytes(), nil
}

func ConstructAppInstanceValues(instance *appv1alpha1.AppInstance) (string, error) {
	var inputs map[string]interface{}
	if err := yaml.Unmarshal([]byte(instance.Spec.Inputs), &inputs); err != nil {
		klog.Errorf("Failed to unmarshal inputs: %v", err)
		return "", err

	}
	inputJson, err := json.Marshal(inputs)
	if err != nil {
		klog.Errorf("Failed to marshal inputs: %v", err)
		return "", err
	}
	return fmt.Sprintf(AppInstanceValues, instance.Name,
		instance.Spec.PublicServiceClass, inputJson), nil
}

func ConstructPublicServiceInstanceValues(instance *servicev1alpha1.PublicServiceInstance) (string, error) {
	var inputs map[string]interface{}
	if err := yaml.Unmarshal([]byte(instance.Spec.Inputs), &inputs); err != nil {
		klog.Errorf("Failed to unmarshal inputs: %v", err)
		return "", err

	}
	inputJson, err := json.Marshal(inputs)
	if err != nil {
		klog.Errorf("Failed to marshal inputs: %v", err)
		return "", err
	}
	return fmt.Sprintf(PublicServiceInstanceValues, instance.Name, inputJson), nil
}

var ifaceNameRegex []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^eth[0-9]+`),
	regexp.MustCompile(`^enps[0-9]+`),
}

func GetLocalServerIPAddress() (string, error) {
	// TBD: This way may not be nice enough, let's find a better way in the future
	ifaces, err := net.Interfaces()
	if err != nil {
		klog.Errorf("Failed to get interfaces: %v", err)
		return "", err
	}
	for _, i := range ifaces {
		for _, r := range ifaceNameRegex {
			if r.MatchString(i.Name) {
				addrs, err := i.Addrs()
				if err != nil {
					klog.Errorf("Failed to get addresses: %v", err)
					return "", err
				}
				for _, addr := range addrs {
					ip, _, err := net.ParseCIDR(addr.String())
					if err != nil {
						klog.Errorf("Failed to parse CIDR: %v", err)
						return "", err
					}
					if ip.To4() != nil {
						return ip.String(), nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("failed to find local server IP address")
}

func CleanInstanceDerivedServiceResource(client kubernetes.Interface, name, namespace string) error {
	svc, err := client.CoreV1().Services(namespace).
		Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("Failed to get service: %v", err)
		return err
	}
	if svc.Finalizers != nil {
		// K8s will add a finalizer to the service, we should delete it first
		svc.Finalizers = nil
		_, err = client.CoreV1().Services(namespace).
			Update(context.Background(), svc, metav1.UpdateOptions{})
		if err != nil {
			klog.Errorf("Failed to delete service's finalizers: %v", err)
			return err
		}
	}

	err = client.CoreV1().Services(namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("Failed to delete service: %v", err)
		return err
	}

	return nil
}

func CleanInstanceDerivedConfigMapResource(client kubernetes.Interface, name, namespace string) error {
	err := client.CoreV1().ConfigMaps(namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("Failed to delete configmap: %v", err)
		return err
	}
	return nil
}

func CleanInstanceDerivedStatefulSetResource(client kubernetes.Interface, name, namespace string) error {
	err := client.AppsV1().StatefulSets(namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("Failed to delete statefulset: %v", err)
		return err
	}
	return nil
}

func CreateOrUpdateService(client kubernetes.Interface,
	manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	labels map[string]string) error {
	var svc apicorev1.Service
	if err := yaml.Unmarshal(manifestContent, &svc); err != nil {
		klog.Errorf("Failed to unmarshal service: %v", err)
		return nil
	}
	svc.Labels = labels
	*derivedResoruce = append(*derivedResoruce, commonv1alpha1.DerivedResource{
		APIVersion: svc.APIVersion,
		Kind:       svc.Kind,
		Name:       svc.Name,
	})

	svcExist, err := client.CoreV1().Services(svc.Namespace).
		Get(context.Background(), svc.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("Failed to get service: %v", err)
			return err
		}
		_, err = client.CoreV1().Services(svc.Namespace).
			Create(context.Background(), &svc, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("Failed to create service: %v", err)
			return err
		}
		return nil
	}
	svcCopy := svcExist.DeepCopy()
	svcCopy.Spec = svc.Spec
	svcCopy.Labels = svc.Labels
	_, err = client.CoreV1().Services(svcCopy.Namespace).
		Update(context.Background(), svcCopy, metav1.UpdateOptions{})
	return err
}

func CreateOrUpdateConfigmap(client kubernetes.Interface,
	manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	labels map[string]string) error {
	var cm apicorev1.ConfigMap
	if err := yaml.Unmarshal(manifestContent, &cm); err != nil {
		klog.Errorf("Failed to unmarshal configmap: %v", err)
		return err
	}
	cm.Labels = labels
	*derivedResoruce = append(*derivedResoruce, commonv1alpha1.DerivedResource{
		APIVersion: cm.APIVersion,
		Kind:       cm.Kind,
		Name:       cm.Name,
	})

	cmExist, err := client.CoreV1().ConfigMaps(cm.Namespace).
		Get(context.Background(), cm.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("Failed to get configmap: %v", err)
			return err
		}
		_, err = client.CoreV1().ConfigMaps(cm.Namespace).
			Create(context.Background(), &cm, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("Failed to create configmap: %v", err)
			return err
		}
		return nil
	}

	cmCopy := cmExist.DeepCopy()
	cmCopy.Data = cm.Data
	cmCopy.Labels = cm.Labels
	_, err = client.CoreV1().ConfigMaps(cmCopy.Namespace).
		Update(context.Background(), cmCopy, metav1.UpdateOptions{})
	return err
}

func CreateOrUpdateStatefulset(client kubernetes.Interface,
	manifestContent []byte,
	derivedResoruce *[]commonv1alpha1.DerivedResource,
	labels map[string]string) error {
	var sts v1.StatefulSet
	if err := yaml.Unmarshal(manifestContent, &sts); err != nil {
		klog.Errorf("Failed to unmarshal statefulset: %v", err)
		return err
	}
	sts.Labels = labels
	*derivedResoruce = append(*derivedResoruce, commonv1alpha1.DerivedResource{
		APIVersion: sts.APIVersion,
		Kind:       sts.Kind,
		Name:       sts.Name,
	})

	stsExist, err := client.AppsV1().StatefulSets(sts.Namespace).
		Get(context.Background(), sts.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("Failed to get statefulset: %v", err)
			return err
		}
		_, err = client.AppsV1().StatefulSets(sts.Namespace).
			Create(context.Background(), &sts, metav1.CreateOptions{})
		if err != nil {
			klog.Errorf("Failed to create statefulset: %v", err)
			return err
		}
		return nil
	}
	if stsExist.Labels != nil &&
		stsExist.Labels[InstanceGenerationLabelKey] == labels[InstanceGenerationLabelKey] {
		return nil
	}

	// In order to make sure the sts will be restart, we should recreate it
	err = client.AppsV1().StatefulSets(sts.Namespace).
		Delete(context.Background(), stsExist.Name, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("Failed to delete statefulset: %v", err)
		return err
	}
	_, err = client.AppsV1().StatefulSets(sts.Namespace).
		Create(context.Background(), &sts, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Failed to create statefulset: %v", err)
		return err
	}
	return nil
}
