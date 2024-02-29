package main

import (
	"os"

	pkgserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"

	"github.com/openapp-dev/openapp/cmd/frp4/app"
)

func main() {
	ctx := pkgserver.SetupSignalContext()
	cmd := app.NewFrp4GatewayCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
