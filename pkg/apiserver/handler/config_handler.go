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
	UserName string `json:"username"`
	Password string `json:"password"`
}

func GetConfigHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get config...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	cfg, err := openappHelper.ConfigMapLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
	if err != nil {
		klog.Errorf("Failed to get config: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get config", nil)
		return
	}

	resp := &OpenAPPSystemConfig{}
	resp.Registry = cfg.Data["registry"]
	resp.UserName = cfg.Data["username"]
	resp.Password = cfg.Data["password"]
	utils.ReturnFormattedData(ctx, http.StatusOK, "Get config successfully", resp)
}

func UpdateConfigHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to update config...")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		klog.Errorf("Failed to read request body: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to read request body", nil)
		return
	}
	config := &OpenAPPSystemConfig{}
	if err := json.Unmarshal(body, config); err != nil {
		klog.Errorf("Failed to unmarshal request body: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to unmarshal request body", nil)
		return
	}

	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}
	systemCfg, err := openappHelper.ConfigMapLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
	if err != nil {
		klog.Errorf("Failed to get config: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get config", nil)
		return
	}
	updatedCfg := systemCfg.DeepCopy()
	updatedCfg.Data["registry"] = config.Registry
	updatedCfg.Data["username"] = config.UserName
	updatedCfg.Data["password"] = config.Password

	if _, err := openappHelper.K8sClient.CoreV1().ConfigMaps(utils.SystemNamespace).Update(context.TODO(),
		updatedCfg, metav1.UpdateOptions{}); err != nil {
		klog.Errorf("Failed to update config: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to update config", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Config updated successfully", nil)
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
