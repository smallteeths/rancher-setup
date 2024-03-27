package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rancher-setup/app/request"
	"github.com/rancher-setup/internal/variable"
	"github.com/rancher-setup/internal/variable/consts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"os/exec"
)

type Index struct {
	base
}

func (i *Index) InstallRancher(ctx *gin.Context) {
	var param request.RancherConfig
	if err := i.base.Validate(ctx, &param); err == nil {
		if ackErr := getAckStatus(); ackErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": ackErr})
			return
		}
		rancherPodExists, _, err := getRancherpod()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if rancherPodExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Rancher already installed"})
			return
		}
		cmd := exec.Command("bash", "-c", consts.Scripts)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("HarborUser=%s", param.HarborUser),
			fmt.Sprintf("HarborPasswd=%s", param.HarborPasswd),
			fmt.Sprintf("Host=%s", param.Host),
			fmt.Sprintf("Version=%s", param.Version))
		output, err := cmd.CombinedOutput()
		if err != nil {
			variable.Log.Info(err.Error())
			variable.Log.Info(string(output))
		}
		variable.Log.Info(string(output))
		variable.Config.Cache("rancherConfig", &param)
		ctx.JSON(http.StatusOK, gin.H{"data": param})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}

func (i *Index) GetRancherState(ctx *gin.Context) {
	rancherPodExists, rancherPodActive, err := getRancherpod()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{
		"rancherPodExists": rancherPodExists,
		"rancherPodActive": rancherPodActive,
	}})
}

func (i *Index) GetRancherConfig(ctx *gin.Context) {
	if config, ok := variable.Config.Get("rancherConfig").(*request.RancherConfig); ok {
		ctx.JSON(http.StatusOK, gin.H{"data": config})
	} else {
		emptyConfig := &request.RancherConfig{}
		ctx.JSON(http.StatusOK, gin.H{"data": emptyConfig})
	}
}

func getRancherpod() (bool, bool, error) {
	clientset, err := variable.Config.GetClientset()
	if err != nil {
		return false, false, err
	}
	pods, err := clientset.CoreV1().Pods("cattle-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list pods: %v\n", err)
		return false, false, err
	}
	rancherPodExists := false
	for _, pod := range pods.Items {
		if _, ok := pod.Labels["app"]; ok && pod.Labels["app"] == "rancher" {
			rancherPodExists = true
			break
		}
	}
	allActive := false
	if rancherPodExists {
		allActive = true
		for _, pod := range pods.Items {
			if _, ok := pod.Labels["app"]; ok && pod.Labels["app"] == "rancher" {
				isActive := pod.Status.Phase == corev1.PodRunning
				if !isActive {
					allActive = false
					break
				}
			}

		}
	}

	return rancherPodExists, allActive, nil
}
func getAckStatus() error {
	clientset, err := variable.Config.GetClientset()
	if err != nil {
		return err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodes.Items {
		if !isNodeReady(node) {
			return errors.New("Ack is not active")
		}
	}
	return nil
}

func isNodeReady(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}
