package handler

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/utils"
)

func GetOpenAPPVersionHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to query openapp version...")

	version := utils.GetOpenAPPVersion()
	returnFormattedData(ctx, 200, "Get openapp version successfully", version)
}
