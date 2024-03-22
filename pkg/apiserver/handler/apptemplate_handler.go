package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/utils"
)

func ListAllAppTemplatesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all app templates...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	appTemps, err := openappHelper.AppTemplateLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list app templates: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to list app templates", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "List app templates successfully", appTemps)
}

func GetAppTemplateHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get app template...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	tempName := ctx.Param("templateName")
	appTemp, err := openappHelper.AppTemplateLister.Get(tempName)
	if err != nil {
		klog.Errorf("Failed to get app template: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get app template", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Get app template successfully", appTemp)
}
