package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"

	appv1alpha1 "github.com/openapp-dev/openapp/pkg/apis/app/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/utils"
)

func ListAllAppInstancesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	appIns, err := openappHelper.AppInstanceLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list app instances: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "List app instances successfully", appIns)
}

func GetAppInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	insName := ctx.Param("instanceName")
	appIns, err := openappHelper.AppInstanceLister.AppInstances(utils.InstanceNamespace).
		Get(insName)
	if err != nil {
		klog.Errorf("Failed to get app instance: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Get app instance successfully", appIns)
}

func DeleteAppInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to delete app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	insName := ctx.Param("instanceName")
	err = openappHelper.OpenAPPClient.AppV1alpha1().
		AppInstances(utils.InstanceNamespace).
		Delete(context.Background(), insName, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("Failed to delete app instance: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Delete app instance successfully", nil)
}

func CreateOrUpdateAppInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to create or update app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var appIns appv1alpha1.AppInstance
	insJsonBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		klog.Errorf("Failed to read request body: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if err := json.Unmarshal(insJsonBody, &appIns); err != nil {
		klog.Errorf("Failed to unmarshal request body: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	appIns.Name = ctx.Param("instanceName")
	appIns.Namespace = utils.InstanceNamespace
	_, err = openappHelper.OpenAPPClient.AppV1alpha1().AppInstances(utils.InstanceNamespace).
		Create(context.Background(), &appIns, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Failed to create or update app instance: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Create or update app instance successfully", nil)
}

func AppInstanceLoggingHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to log app instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	insName := ctx.Param("instanceName")
	pods, err := openappHelper.K8sClient.CoreV1().Pods(utils.InstanceNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + insName,
	})
	if err != nil {
		klog.Errorf("Failed to get app instance's pod: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if len(pods.Items) == 0 || len(pods.Items) > 1 {
		klog.Warningf("No resource found for app instance(%s)", insName)
		utils.ReturnFormattedData(ctx, http.StatusNotFound, "No resource found for app instance", nil)
		return
	}

	req := openappHelper.K8sClient.CoreV1().Pods(utils.InstanceNamespace).GetLogs(pods.Items[0].Name, &v1.PodLogOptions{
		Container: pods.Items[0].Spec.Containers[0].Name,
	})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		klog.Errorf("Failed to get app instance's pod logs: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer podLogs.Close()

	streamUsed := ctx.Query("stream") == "true"
	if streamUsed {
		ctx.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, podLogs)
			if err != nil {
				klog.Errorf("Error streaming logs: %v", err)
				return false
			}
			return true
		})
	} else {
		logs, err := io.ReadAll(podLogs)
		if err != nil {
			klog.Errorf("Error reading logs: %v", err)
			utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utils.ReturnFormattedData(ctx, http.StatusOK, "Get app instance logs successfully", string(logs))
	}
}
