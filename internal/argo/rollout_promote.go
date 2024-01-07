package argo

import (
	"context"
	"log"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PromoteRollout(argoClientset *argoclientset.Clientset, rolloutObj *v1alpha1.Rollout, namespace string) error {
	AutoPromotionEnabled := true
	rolloutObj.Spec.Strategy.BlueGreen.AutoPromotionEnabled = &AutoPromotionEnabled
	_, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Update(context.TODO(), rolloutObj, metav1.UpdateOptions{})
	if err != nil {
		log.Println("Error when promoting the Argo Rollout:", err)
		return err
	}
	return nil
}

func GetPromoteStatus(rolloutObj *v1alpha1.Rollout) (bool, string) {
	log.Printf("Phase: %s", rolloutObj.Status.Phase)
	if rolloutObj.Status.Phase == "Healthy" {
		return false, "Failed, pre-prod is not ready in production yet, please check your build status!"
	} else if rolloutObj.Status.Phase == "Progressing" {
		return false, "Hang tight, pre-prod deployment is still going! if this happens for more than 3 mins, means preprod got crashed"
	} else if rolloutObj.Status.Phase == "Paused" {
		return true, "Ready to promote"
	}

	return false, "Unknown pre-prod status, please report it"
}
