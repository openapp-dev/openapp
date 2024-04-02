package app

import (
	"context"
	"flag"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/controller/appinstance"
	"github.com/openapp-dev/openapp/pkg/controller/apptemplate"
	"github.com/openapp-dev/openapp/pkg/controller/publicserviceinstance"
	"github.com/openapp-dev/openapp/pkg/controller/publicservicetemplate"
	"github.com/openapp-dev/openapp/pkg/controller/registry"
	"github.com/openapp-dev/openapp/pkg/controller/types"
	openappclient "github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

func NewControllerCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "openapp-controller",
		Long: `openapp-controller used to control the openapp resources`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(ctx); err != nil {
				return err
			}
			return nil
		},
	}

	fss := cliflag.NamedFlagSets{}
	logFlagSet := fss.FlagSet("log")
	klog.InitFlags(flag.CommandLine)
	logFlagSet.AddGoFlagSet(flag.CommandLine)
	cmd.Flags().AddFlagSet(logFlagSet)

	return cmd
}

var controllerNewFuncList = []types.NewControllerFunc{
	apptemplate.NewAppTempalteController,
	publicservicetemplate.NewPublicServiceTemplateController,
	appinstance.NewAppInstanceController,
	appinstance.NewAppInstanceServiceController,
	appinstance.NewAppInstanceStatusController,
	publicserviceinstance.NewPublicServiceInstanceController,
	publicserviceinstance.NewPublicServiceInstanceStatusController,
	publicserviceinstance.NewPublicServiceInstanceServiceController,
	registry.NewRegistryController,
}

func run(ctx context.Context) error {
	version := utils.GetOpenAPPVersion()
	klog.Infof("Start openapp-controller, version: %v, commit: %s", version.GitVersion, version.GitCommit)
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Failed to get in-cluster config: %v", err)
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create k8s client: %v", err)
	}
	openappClient, err := openappclient.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create client: %v", err)
	}
	openappHelper := utils.NewOpenAPPHelper(ctx, k8sClient, openappClient)

	ctls := []types.ControllerInterface{}
	for _, controllerFunc := range controllerNewFuncList {
		ctls = append(ctls, controllerFunc(openappHelper))
	}

	klog.Infof("Wait resource cache sync...")
	if ok := cache.WaitForCacheSync(ctx.Done(),
		openappHelper.ConfigMapInformer.HasSynced,
		openappHelper.AppInstanceInformer.HasSynced,
		openappHelper.PublicServiceInstanceInformer.HasSynced,
		openappHelper.ServiceInformer.HasSynced,
		openappHelper.StatefulSetInformer.HasSynced); !ok {
		klog.Fatal("Failed to wait for cache sync")
	}

	for _, ctl := range ctls {
		go ctl.Start()
	}

	<-ctx.Done()
	return nil
}
