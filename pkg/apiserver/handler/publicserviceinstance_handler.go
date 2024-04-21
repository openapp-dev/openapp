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

	servicev1alpha1 "github.com/openapp-dev/openapp/pkg/apis/service/v1alpha1"
	"github.com/openapp-dev/openapp/pkg/utils"
)

func ListAllPublicServiceInstancesHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to list all public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	publicServiceIns, err := openappHelper.PublicServiceInstanceLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list public service instances: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to list public service instances", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "List public service instances successfully", publicServiceIns)
}

func GetPublicServiceInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to get public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	insName := ctx.Param("instanceName")
	ins, err := openappHelper.PublicServiceInstanceLister.PublicServiceInstances(utils.InstanceNamespace).
		Get(insName)
	if err != nil {
		klog.Errorf("Failed to get publicservice instance:%v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get publicservice instance", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Get publicservice instance successfully", ins)
}

func DeletePublicServiceInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to delete public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to get openapp lister", nil)
		return
	}

	insName := ctx.Param("instanceName")
	err = openappHelper.OpenAPPClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).
		Delete(context.Background(), insName, metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("Failed to delete public service instance: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to delete public service instance", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Delete public service instance successfully", nil)
}

func CreateOrUpdatePublicServiceInstanceHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to create or update public service instance...")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		klog.Errorf("Failed to read request body: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to read request body", nil)
		return
	}

	var ins servicev1alpha1.PublicServiceInstance
	if err := json.Unmarshal(body, &ins); err != nil {
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

	ins.Name = ctx.Param("instanceName")
	ins.Namespace = utils.InstanceNamespace
	_, err = openappHelper.OpenAPPClient.ServiceV1alpha1().PublicServiceInstances(utils.InstanceNamespace).
		Create(context.Background(), &ins, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Failed to create public service instance: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, "Failed to create public service instance", nil)
		return
	}

	utils.ReturnFormattedData(ctx, http.StatusOK, "Create public service instance successfully", nil)
}

func PublicServiceInstanceLoggingHandler(ctx *gin.Context) {
	klog.V(4).Infof("Start to log public service instance...")
	openappHelper, err := getOpenAPPHelper(ctx)
	if err != nil {
		klog.Errorf("Failed to get openapp lister: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	insName := ctx.Param("instanceName")
	pods, err := openappHelper.K8sClient.CoreV1().Pods(utils.InstanceNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "publicservice=" + insName,
	})
	if err != nil {
		klog.Errorf("Failed to get public service instance's pod: %v", err)
		utils.ReturnFormattedData(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if len(pods.Items) == 0 || len(pods.Items) > 1 {
		klog.Warningf("No resource found for public service instance(%s)", insName)
		utils.ReturnFormattedData(ctx, http.StatusNotFound, "No resource found for public service instance", nil)
		return
	}

	// TBD: Maybe we can let user choose
	req := openappHelper.K8sClient.CoreV1().Pods(utils.InstanceNamespace).GetLogs(pods.Items[0].Name, &v1.PodLogOptions{
		Container: pods.Items[0].Spec.Containers[0].Name,
	})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		klog.Errorf("Failed to get public service instance's pod logs: %v", err)
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
		utils.ReturnFormattedData(ctx, http.StatusOK, "Get public service instance logs successfully", string(logs))
	}
}
