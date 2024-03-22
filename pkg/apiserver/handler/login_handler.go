package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"

	"github.com/openapp-dev/openapp/pkg/utils"
)

type loginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(ctx *gin.Context) {
	klog.V(4).Info("Start to login...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}
	json := make(map[string]interface{})
	err = ctx.BindJSON(&json)
	if err != nil {
		klog.Errorf("Failed to bind json: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusBadRequest, "Failed to bind json", nil)
		return
	}
	username, ok := json["username"].(string)
	if !ok {
		klog.Errorf("Failed to get username from json")
		utils.ReturnFormattedData(ctx, http.StatusBadRequest, "Failed to get username from json", nil)
		return
	}
	password, ok := json["password"].(string)
	if !ok {
		klog.Errorf("Failed to get password from json")
		utils.ReturnFormattedData(ctx, http.StatusBadRequest, "Failed to get password from json", nil)
		return
	}
	cfg, err := openappHelper.ConfigMapLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
	if err != nil {
		klog.Errorf("Failed to get config: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get config", nil)
		return
	}
	if cfg.Data["username"] != username || cfg.Data["password"] != password {
		klog.Errorf("Failed to login")
		utils.ReturnFormattedData(ctx, http.StatusUnauthorized, "Failed to login", nil)
		return
	}
	token, err := utils.NewJWT([]byte(cfg.Data["password"])).GenerateToken(username, password)
	if err != nil {
		klog.Errorf("Failed to generate token: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}
	utils.ReturnFormattedData(ctx, http.StatusOK, "Login successfully", &loginResponse{Token: token})
}
