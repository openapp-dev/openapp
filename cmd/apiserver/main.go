package main

import (
	"os"

	pkgserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"

	"github.com/openapp-dev/openapp/cmd/apiserver/app"
)

func main() {
	ctx := pkgserver.SetupSignalContext()
	cmd := app.NewApiServerCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
