package main

import (
	"fmt"
	"log"

	"cubeflow/internal/api"
	"cubeflow/internal/kubernetes"
	"cubeflow/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	err := config.Load("./config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	route := gin.Default()
	internalApi := route.Group("/v1")
	externalApi := route.Group("/v2")

	argoClientConfigs, kubeClientConfigs, nil := kubernetes.GetKubeConfigClient()
	api.InternalAPIRoutes(internalApi, kubeClientConfigs, argoClientConfigs)
	api.ExternalAPIRoutes(externalApi, kubeClientConfigs, argoClientConfigs)

	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)
	log.Printf("Starting CubeFlow service on port %d...", port)
	if err := route.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start CubeFlow service: %v", err)
	}
}
