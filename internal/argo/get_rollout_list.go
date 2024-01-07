package argo

import (
	"context"
	"fmt"
	"log"

	helperFunc "cubeflow/internal/helpers"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RolloutPhase string

func (phase RolloutPhase) String() string {
	return string(phase)
}

func GetRolloutStatusList(rolloutList *v1alpha1.RolloutList) map[string]string {
	keyValueMap := make(map[string]string)
	for _, rollout := range rolloutList.Items {
		log.Printf("RolloutName: %s, status %s", rollout.Name, rollout.Status.Phase)
		phase := RolloutPhase(rollout.Status.Phase)
		keyValueMap[rollout.Name] = phase.String()
	}
	log.Println("Rolout Status Array: ", keyValueMap)
	return keyValueMap
}

func GetRolloutListParsed(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, namespace string) string {
	rolloutList, _ := GetRolloutList(clientset, argoClientset, namespace)
	allRolloutname := GetAllRolloutName(rolloutList)
	strRolloutName := helperFunc.ParseToSingleLine(allRolloutname, ", ")
	return strRolloutName
}

func GetRolloutList(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, namespace string) (*v1alpha1.RolloutList, error) {
	rolloutList, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error when Retrieve the Argo Rollout Get List: ", err)
		return nil, err
	}
	return rolloutList, nil
}

func GetAllRolloutName(rolloutList *v1alpha1.RolloutList) []string {
	var allRolloutName []string
	for _, rollout := range rolloutList.Items {
		allRolloutName = append(allRolloutName, rollout.Name)
	}
	log.Println("Rolout Status Array: ", allRolloutName)
	return allRolloutName
}

func ParseSingleLineArr(arr map[string]string) string {
	keyValueString := ""
	for key, value := range arr {
		if value == "Healthy" {
			value = "NOT READY"
		} else if value == "Paused" {
			value = "READY TO PROMOTE"
		} else if value == "Progressing" {
			value = "PROGRESSING"
		}
		keyValueString += fmt.Sprintf("(%s: %s), ", key, value)
	}
	return keyValueString
}

func GetRolloutManifest(argoClientset *argoclientset.Clientset, namespace string, rolloutName string) (*v1alpha1.Rollout, error) {
	var rollout *v1alpha1.Rollout
	rollout, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Get(context.TODO(), rolloutName, metav1.GetOptions{})
	if err != nil {
		log.Println("Error when Retrieve the Argo Rollout Manifest: ", err)
		return nil, err
	}
	return rollout, nil
}

func GetImage(rolloutManifest *v1alpha1.Rollout) string {
	return rolloutManifest.Spec.Template.Spec.Containers[0].Image
}
