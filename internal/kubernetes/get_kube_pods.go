package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func IsPodsCrashed(clientset *kubernetes.Clientset, namespace string, rolloutName string) (bool, []string, error) {
	var nulArr []string
	podsManifestInNamespace, err := GetPodsManifestInNamespace(clientset, namespace)
	if err != nil {
		log.Fatalln("Can list the pods list: ", err)
		return false, nulArr, err
	}

	podStatus, podCrashedList, err := GetPodsStatus(podsManifestInNamespace, rolloutName)
	if err != nil {
		log.Fatalln("Can list the pods list: ", err)
		return false, podCrashedList, err
	}
	return podStatus, podCrashedList, nil
}

func GetPodsManifestInNamespace(clientset *kubernetes.Clientset, namespace string) (*v1.PodList, error) {
	podsManifest, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Can list the pods list: ", err)
		return nil, err
	}
	return podsManifest, err
}

func GetPodsStatus(pods *v1.PodList, rolloutName string) (bool, []string, error) {
	var crashedPods []string
	oneOrMorePodCrashed := false

	for _, pod := range pods.Items {
		lengthSidecar := len(pod.Spec.Containers)
		for i := 0; i < lengthSidecar; i++ {
			if ContainerNotReady(pod, rolloutName, i) {
				log.Printf("Caught container that not ready, pods: " + pod.Name)
				oneOrMorePodCrashed = true
				crashedPods = append(crashedPods, pod.Name)
				break
			}
		}
	}
	return oneOrMorePodCrashed, crashedPods, nil
}

func ContainerNotReady(pod v1.Pod, rolloutName string, index int) bool {
	return pod.Status.ContainerStatuses[index].Name == rolloutName && !pod.Status.ContainerStatuses[index].Ready
}
