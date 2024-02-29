package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"
)

func ListAllAppTemplatesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all app templates...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	appTemps, err := openappHelper.AppTemplateLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list app templates: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list app templates"})
		return
	}

	ctx.AsciiJSON(http.StatusOK, appTemps)
}

func GetAppTemplateHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get app template...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	tempName := ctx.Param("templateName")
	appTemp, err := openappHelper.AppTemplateLister.Get(tempName)
	if err != nil {
		klog.Errorf("Failed to get app template: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get app template"})
		return
	}

	ctx.AsciiJSON(http.StatusOK, appTemp)
}
