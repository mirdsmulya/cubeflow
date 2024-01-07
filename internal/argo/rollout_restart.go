package argo

import (
	"context"
	"log"
	"net/url"
	"time"

	helperFunc "cubeflow/internal/helpers"
	kubeFunc "cubeflow/internal/kubernetes"

	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func MultipleRestartPods(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, namespace string, prompInput string) (string, string) {
	var errorList []string
	var proceedList []string
	joinText := ", "
	splitParameter := ","
	decodedPrompInput, _ := url.QueryUnescape(prompInput)

	serviceNameList := helperFunc.MutipleExec(decodedPrompInput, splitParameter)
	for _, service := range serviceNameList {
		err := RestartArgoRollout(clientset, argoClientset, namespace, service)
		if err != nil {
			errorList = append(errorList, service)
		} else {
			proceedList = append(proceedList, service)
		}
	}

	errListParsed := helperFunc.ParseToSingleLine(errorList, joinText)
	proceedListParsed := helperFunc.ParseToSingleLine(proceedList, joinText)
	return proceedListParsed, errListParsed
}

func RestartArgoRollout(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, namespace, rolloutName string) error {
	podCrashed, podCrashedList, err := kubeFunc.IsPodsCrashed(clientset, namespace, rolloutName)
	if err != nil {
		log.Println("Restart error when check pods status", err)
		return err
	}

	if podCrashed {
		err := HardRestart(clientset, namespace, podCrashedList)
		if err != nil {
			return err
		}
		log.Printf("Hard restart for %s in namespace %s\n", rolloutName, namespace)
		return nil
	}

	err = NormalRestart(argoClientset, namespace, rolloutName)
	if err != nil {
		return err
	}

	log.Printf("Restarting Argo Rollout %s in namespace %s\n", rolloutName, namespace)
	return nil
}

func HardRestart(clientset *kubernetes.Clientset, namespace string, podCrashedList []string) error {
	for _, pod := range podCrashedList {
		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod, metav1.DeleteOptions{})
		if err != nil {
			log.Println("Restart error when perform delete pods", err)
			return err
		}
	}
	return nil
}

func NormalRestart(argoClientset *argoclientset.Clientset, namespace string, rolloutName string) error {
	rolloutObj, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Get(context.TODO(), rolloutName, metav1.GetOptions{})
	if err != nil {
		log.Println("Error when Retrieve the Argo Rollout object, the prompt might be wrong", err)
		return err
	}

	currentTime := metav1.NewTime(time.Now())
	rolloutObj.Spec.RestartAt = &currentTime
	_, err = argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Update(context.TODO(), rolloutObj, metav1.UpdateOptions{})
	if err != nil {
		log.Println("Error when Apply the changes to the Argo Rollout")
		return err
	}

	return nil
}
