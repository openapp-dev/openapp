package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"
)

func ListAllPublicServiceTemplatesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all publicservice templates...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		returnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	publicServiceTemps, err := openappHelper.PublicServiceTemplateLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list publicservice templates: %v", err)
		returnFormattedData(ctx, http.StatusInternalServerError, "Failed to list publicservice templates", nil)
		return
	}

	returnFormattedData(ctx, http.StatusOK, "List publicservice templates successfully", publicServiceTemps)
}

func GetPublicServiceTemplateHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get publicservice template...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		returnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	tempName := ctx.Param("templateName")
	publicServiceTemp, err := openappHelper.PublicServiceTemplateLister.Get(tempName)
	if err != nil {
		klog.Errorf("Failed to get publicservice template: %v", err)
		returnFormattedData(ctx, http.StatusInternalServerError, "Failed to get publicservice template", nil)
		return
	}

	returnFormattedData(ctx, http.StatusOK, "Get publicservice template successfully", publicServiceTemp)
}
