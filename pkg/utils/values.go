package utils

const (
	OpenAPPHelperKey = "openappHelper"

	RegistryKey                   = "registry"
	RegistryCachePath             = "/root/openapp/registry"
	AppTemplatePath               = "app-template"
	AppTemplateBasePath           = "app-template"
	PublicServiceTemplatePath     = "publicservice-template"
	PublicServiceTemplateBasePath = "publicservice-template"
	TemplateFileName              = "template.yaml"
	TemplateResourceDirName       = "resource"
	TemplateManifestsDirName      = "manifests"

	RegistryUpdateTimeAnnotationKey = "registry.openapp.dev/update-time"
	ServiceExposeClassLabelKey      = "service.openapp.dev/expose-class"
	AppInstanceLabelKey             = "app.openapp.dev/app-instance"
	PublicServiceInstanceLabelKey   = "service.openapp.dev/publicservice-instance"
	InstanceGenerationLabelKey      = "instance.openapp.dev/instance-generation"

	InstanceNamespace = "openapp"
	SystemNamespace   = "openapp-system"
	SystemConfigMap   = "openapp-config"
	VolumeConfigMap   = "volume-config"

	TemplateManifestServiceFile     = "service.yaml"
	TemplateManifestStatefulSetFile = "statefulset.yaml"
	TemplateManifestConfigMapFile   = "configmap.yaml"

	InstanceDerivedResourceServiceKind     = "Service"
	InstanceDerivedResourceStatefulSetKind = "StatefulSet"
	InstanceDerivedResourceConfigMapKind   = "ConfigMap"

	OpenAPPDNSName = "openapp"

	PublicServiceInstanceControllerFinalizerKey = "publicservice-instance-controller"
	AppInstanceControllerFinalizerKey           = "app-instance-controller"
)

var (
	AppInstanceValues = `
	{
		"openapp": {
			"instance_name": "%s",
			"service_class": "%s"
		},
		"inputs": %s
	}`
	PublicServiceInstanceValues = `
	{
		"openapp": {
			"instance_name": "%s"
		},
		"inputs": %s
	}`
)
