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

	servicev1alpha1 "github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
)

func ListAllPublicServiceInstancesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	publicServiceIns, err := openappHelper.PublicServiceInstanceLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list public service instances: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list public service instances"})
		return
	}

	ctx.AsciiJSON(http.StatusOK, publicServiceIns)
}

func GetPublicServiceInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	insName := ctx.Param("instanceName")
	ins, err := openappHelper.PublicServiceInstanceLister.PublicServiceInstances(utils.InstanceNamespace).
		Get(insName)
	if err != nil {
		klog.Errorf("Failed to get publicservice instance:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get publicservice instance"})
		return
	}

	ctx.AsciiJSON(http.StatusOK, ins)
}

func DeletePublicServiceInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to delete public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	insName := ctx.Param("instanceName")
	err = openappHelper.OpenAPPClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).
		Delete(context.Background(), insName, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("Failed to delete public service instance: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete public service instance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "success"})
}

func CreateOrUpdatePublicServiceInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to create or update public service instance...")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		klog.Errorf("Failed to read request body: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	var ins servicev1alpha1.PublicServiceInstance
	if err := json.Unmarshal(body, &ins); err != nil {
		klog.Errorf("Failed to unmarshal request body: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal request body"})
		return
	}

	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	_, err = openappHelper.OpenAPPClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).
		Create(context.Background(), &ins, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Failed to create public service instance: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create public service instance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "success"})
}
