package main

import (
	"os"

	pkgserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"

	"github.com/openapp-dev/openapp/cmd/controller/app"
)

func main() {
	ctx := pkgserver.SetupSignalContext()
	cmd := app.NewControllerCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
