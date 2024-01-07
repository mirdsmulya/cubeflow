package gke

import (
	"context"
	"cubeflow/pkg/config"
	"log"

	container "google.golang.org/api/container/v1beta1"
	"google.golang.org/api/option"
)

func ScaleNodeGroup(gkeNodePoolName string, targetNodeCount int64) error {

	ctx := context.Background()

	// Authenticate with Google Cloud using the service account key
	client, err := container.NewService(ctx, option.WithCredentialsFile(config.Variable.GCP.CredentialsPath))
	if err != nil {
		log.Printf("Failed to create GKE client: %v", err)
	}

	// Create a SetNodePoolSizeRequest
	setNodePoolSizeRequest := &container.SetNodePoolSizeRequest{
		NodeCount: int64(targetNodeCount),
	}

	// Update the node pool configuration
	_, err = client.Projects.Zones.Clusters.NodePools.SetSize(config.Variable.GCP.ProjectID, config.Variable.GCP.GKE.Zone, config.Variable.GCP.GKE.ClusterName, gkeNodePoolName, setNodePoolSizeRequest).Context(ctx).Do()
	if err != nil {
		log.Printf("Error updating node pool configuration: %v", err)
		return err
	}
	log.Printf("successfully scale nodepool to : %v", targetNodeCount)
	return nil
}
