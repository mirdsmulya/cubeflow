package kubernetes

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ScaleDownAllStatefulsetsNamespaces(clientset *kubernetes.Clientset, namespaces []string) error {

	// Iterate over each namespace
	for _, namespace := range namespaces {
		// Process each namespace, for example, print it
		log.Printf("Processing namespace: %v", namespace)

		err := ScaleDownAllStatefulsets(clientset, namespace)
		if err != nil {
			log.Printf("Error scaling down Statefulset: %v", err)
		}
	}

	log.Println("Scaled down all Statefulset completed")

	return nil
}

func ScaleDownAllStatefulsets(clientset *kubernetes.Clientset, namespace string) error {
	statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Iterate through each Statefulset and scale down if replicas > 1
	for _, statefulSet := range statefulSets.Items {
		if *statefulSet.Spec.Replicas > 1 {

			log.Printf("Scaling down statefulset : %v", statefulSet.Name)

			err := ScaleDownStatefulSet(clientset, namespace, statefulSet.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ScaleDownStatefulSet(clientset *kubernetes.Clientset, namespace, statefulSetName string) error {
	// Retrieve the current Statefulset
	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Modify the replicas field
	replicas := int32(1) // Set the desired number of replicas
	statefulSet.Spec.Replicas = &replicas

	// Update the Statefulset
	_, err = clientset.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// Wait for the Statefulset to scale down
	err = WaitForstatefulSetScale(clientset, namespace, statefulSetName, replicas)
	if err != nil {
		return err
	}

	return nil
}

func WaitForstatefulSetScale(clientset *kubernetes.Clientset, namespace, statefulSetName string, replicas int32) error {
	// Timeout for waiting
	timeout := 5 * time.Minute
	deadline := time.Now().Add(timeout)

	for {
		// Retrieve the current Statefulset
		statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Check if the replicas have scaled down
		if *statefulSet.Spec.Replicas == replicas {
			return nil
		}

		// Check timeout
		if time.Now().After(deadline) {
			log.Printf("timeout waiting for Statefulset %v in namespace %v to scale down", statefulSetName, namespace)
		}

		// Sleep for a short interval before checking again
		time.Sleep(5 * time.Second)
	}
}
