package argo

import (
	"context"
	"log"
	"time"

	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ScaleDownAllRolloutsNamespaces(argoClientset *argoclientset.Clientset, namespaces []string) error {

	// Iterate over each namespace
	for _, namespace := range namespaces {
		// Process each namespace, for example, print it
		log.Printf("Processing namespace: %v", namespace)

		err := ScaleDownAllRollouts(argoClientset, namespace)
		if err != nil {
			log.Printf("Error scaling down Rollout: %v", err)
		}
	}

	log.Println("Scaled down all Rollout completed")

	return nil
}

func ScaleDownAllRollouts(argoClientset *argoclientset.Clientset, namespace string) error {
	rollouts, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Iterate through each Rollouts and scale down if replicas > 1
	for _, rollout := range rollouts.Items {
		if *rollout.Spec.Replicas > 1 {
			err := ScaleDownRollout(argoClientset, namespace, rollout.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ScaleDownRollout(argoClientset *argoclientset.Clientset, namespace, rolloutName string) error {
	// Retrieve the current Rollout
	rollout, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Get(context.TODO(), rolloutName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Modify the replicas field
	replicas := int32(1) // Set the desired number of replicas
	rollout.Spec.Replicas = &replicas

	// Update the Rollout
	_, err = argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Update(context.TODO(), rollout, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// Wait for the Rollout to scale down
	err = WaitForRolloutScale(argoClientset, namespace, rolloutName, replicas)
	if err != nil {
		return err
	}

	return nil
}

func WaitForRolloutScale(argoClientset *argoclientset.Clientset, namespace, rolloutName string, replicas int32) error {
	// Timeout for waiting
	timeout := 5 * time.Minute
	deadline := time.Now().Add(timeout)

	for {
		// Retrieve the current Rollout
		rollout, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Get(context.TODO(), rolloutName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Check if the replicas have scaled down
		if *rollout.Spec.Replicas == replicas {
			return nil
		}

		// Check timeout
		if time.Now().After(deadline) {
			log.Printf("timeout waiting for Rollout %v in namespace %v to scale down", rolloutName, namespace)
		}

		// Sleep for a short interval before checking again
		time.Sleep(5 * time.Second)
	}
}
