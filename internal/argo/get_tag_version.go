package argo

import (
	"strings"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
)

func GetTagVersion(argoClientset *argoclientset.Clientset, namespace string, rolloutName string) (string, error) {
	rolloutManifest, err := GetRolloutManifest(argoClientset, namespace, rolloutName)
	if err != nil {
		return "", err
	}
	containerIndex := ContainerIndexCheck(rolloutManifest, rolloutName)
	imageTag := GetTag(rolloutManifest, containerIndex)
	return imageTag, nil
}

func ContainerIndexCheck(rolloutManifest *v1alpha1.Rollout, rolloutName string) int {
	containers := rolloutManifest.Spec.Template.Spec.Containers
	for index, pod := range containers {
		if pod.Name == rolloutName {
			return index
		}
	}
	return 10
}

func GetTag(rolloutManifest *v1alpha1.Rollout, index int) string {
	imagePath := rolloutManifest.Spec.Template.Spec.Containers[index].Image
	tagArr := strings.Split(imagePath, ":")
	tag := tagArr[1]
	return tag
}
