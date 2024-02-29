package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openapp-dev/openapp/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"

	appv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
)

func ListAllAppInstancesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	appIns, err := openappHelper.AppInstanceLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list app instances: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list app instances"})
		return
	}

	ctx.AsciiJSON(http.StatusOK, appIns)
}

func GetAppInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	insName := ctx.Param("instanceName")
	appIns, err := openappHelper.AppInstanceLister.AppInstances(utils.InstanceNamespace).
		Get(insName)
	if err != nil {
		klog.Errorf("Failed to get app instance: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get app instance"})
		return
	}

	ctx.AsciiJSON(http.StatusOK, appIns)
}

func DeleteAppInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to delete app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	insName := ctx.Param("instanceName")
	err = openappHelper.OpenAPPClient.AppV1alpha1().
		AppInstances(utils.InstanceNamespace).
		Delete(context.Background(), insName, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("Failed to delete app instance: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete app instance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Delete app instance successfully"})
}

func CreateOrUpdateAppInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to create or update app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	var appIns appv1alpha1.AppInstance
	insJsonBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		klog.Errorf("Failed to read request body: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	if err := json.Unmarshal(insJsonBody, &appIns); err != nil {
		klog.Errorf("Failed to unmarshal request body: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal request body"})
		return
	}
	appIns.Name = ctx.Param("instanceName")
	appIns.Namespace = utils.InstanceNamespace
	_, err = openappHelper.OpenAPPClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).
		Create(context.Background(), &appIns, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Failed to create or update app instance: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update app instance"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Create or update app instance successfully"})
}
