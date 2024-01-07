package kubernetes

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ScaleDownAllDeploymentsNamespaces(clientset *kubernetes.Clientset, namespaces []string) error {

	// Iterate over each namespace
	for _, namespace := range namespaces {
		// Process each namespace, for example, print it
		log.Printf("Processing namespace: %v", namespace)

		err := ScaleDownAllDeployments(clientset, namespace)
		if err != nil {
			log.Printf("Error scaling down deployment: %v", err)
		}
	}

	log.Println("Scaled down all Deployment completed")

	return nil
}

func ScaleDownAllDeployments(clientset *kubernetes.Clientset, namespace string) error {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Iterate through each Deployment and scale down if replicas > 1
	for _, deployment := range deployments.Items {
		if *deployment.Spec.Replicas > 1 {

			log.Printf("Scaling down statefulset : %v", deployment.Name)

			err := ScaleDownDeployment(clientset, namespace, deployment.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ScaleDownDeployment(clientset *kubernetes.Clientset, namespace, deploymentName string) error {
	// Retrieve the current Deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Modify the replicas field
	replicas := int32(1) // Set the desired number of replicas
	deployment.Spec.Replicas = &replicas

	// Update the Deployment
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// Wait for the Deployment to scale down
	err = WaitForDeploymentScale(clientset, namespace, deploymentName, replicas)
	if err != nil {
		return err
	}

	return nil
}

func WaitForDeploymentScale(clientset *kubernetes.Clientset, namespace, deploymentName string, replicas int32) error {
	// Timeout for waiting
	timeout := 5 * time.Minute
	deadline := time.Now().Add(timeout)

	for {
		// Retrieve the current Deployment
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Check if the replicas have scaled down
		if *deployment.Spec.Replicas == replicas {
			return nil
		}

		// Check timeout
		if time.Now().After(deadline) {
			log.Printf("timeout waiting for Deployment %v in namespace %v to scale down", deploymentName, namespace)
		}

		// Sleep for a short interval before checking again
		time.Sleep(5 * time.Second)
	}
}
