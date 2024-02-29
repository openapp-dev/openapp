package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/utils"
)

type OpenAPPSystemConfig struct {
	Registry string `json:"registry"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func GetConfigHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get config...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get openapp lister"})
		return
	}

	cfg, err := openappHelper.ConfigMapLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
	if err != nil {
		klog.Errorf("Failed to get config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}

	resp := &OpenAPPSystemConfig{}
	resp.Registry = cfg.Data["registry"]
	resp.UserName = cfg.Data["userName"]
	resp.Password = cfg.Data["password"]
	ctx.AsciiJSON(http.StatusOK, resp)
}

func UpdateConfigHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to update config...")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		klog.Errorf("Failed to read request body: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	config := &OpenAPPSystemConfig{}
	if err := json.Unmarshal(body, config); err != nil {
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
	systemCfg, err := openappHelper.ConfigMapLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
	if err != nil {
		klog.Errorf("Failed to get config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}
	updatedCfg := systemCfg.DeepCopy()
	updatedCfg.Data["registry"] = config.Registry
	updatedCfg.Data["userName"] = config.UserName
	updatedCfg.Data["password"] = config.Password

	if _, err := openappHelper.K8sClient.CoreV1().ConfigMaps(utils.SystemNamespace).Update(context.TODO(),
		updatedCfg, metav1.UpdateOptions{}); err != nil {
		klog.Errorf("Failed to update config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Config updated successfully"})
}

func getOpenAPPHelper(ctx *gin.Context) (*utils.OpenAPPHelper, error) {
	lister, ok := ctx.Get(utils.OpenAPPHelperKey)
	if !ok {
		return nil, fmt.Errorf("failed to get openapp lister from context")
	}
	openappLister, ok := lister.(*utils.OpenAPPHelper)
	if !ok {
		return nil, fmt.Errorf("failed to convert openapp lister from context")
	}
	return openappLister, nil
}
