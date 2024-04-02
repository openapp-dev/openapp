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

	"github.com/openapp-dev/openapp/pkg/apiserver/router"
	openappclient "github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

func NewApiServerCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "openapp-apiserver",
		Long: `openapp-apiserver used to serve as the backend`,
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

func run(ctx context.Context) error {
	version := utils.GetOpenAPPVersion()
	klog.Infof("Start openapp-apiserver, version: %s, commit: %s", version.GitVersion, version.GitCommit)
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
	klog.Infof("Wait resource cache sync...")
	if ok := cache.WaitForCacheSync(ctx.Done(),
		openappHelper.ConfigMapInformer.HasSynced,
		openappHelper.AppInstanceInformer.HasSynced,
		openappHelper.PublicServiceInstanceInformer.HasSynced); !ok {
		klog.Fatal("Failed to wait for cache sync")
	}

	openappRouter := router.NewOpenAPPServerRouter(k8sClient, openappClient, openappHelper)
	if err := openappRouter.Run(":8080"); err != nil {
		klog.Fatalf("Run openapp router failed: %v", err)
	}

	<-ctx.Done()
	return nil
}
