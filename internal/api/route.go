package api

import (
	"net/http"

	"cubeflow/internal/handlers"
	"cubeflow/pkg/config"

	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

func InternalAPIRoutes(api *gin.RouterGroup, clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset) {
	api.POST("/service/restart/:namespace/:rolloutName", handlers.RestartPodsV1(clientset, argoClientset))
	api.POST("/db/backup/:dbname", handlers.DatabaseBackupV1())
	api.POST("/cluster/scale/:nodeGroupName/:nodeCount", handlers.ScaleNodeGroupV1())
	api.POST("/service/scaledown", handlers.ScaleDownServiceV1(clientset, argoClientset, config.Variable.GCP.GKE.NamespaceToScale))
	api.POST("/service/getpods/:namespace", handlers.GetPods(clientset))
}

func ExternalAPIRoutes(api *gin.RouterGroup, clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset) {
	api.POST("/service/restart/:namespace", handlers.RestartPodsV2(clientset, argoClientset, config.Variable.Environment))
	api.POST("/service/promote/:namespace", handlers.PromoteRollout(clientset, argoClientset, config.Variable.Environment))
	api.POST("/service/sync", handlers.SyncArgoApp(clientset, argoClientset, config.Variable.Environment))
	api.POST("/db/backup/:dbname", handlers.DatabaseBackupV2(config.Variable.Environment))
	api.POST("/service/rollback/:namespace", handlers.RollbackService(argoClientset, config.Variable.Environment))
	api.GET("/ping", Ping())
}

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		htmlResponse := "<html><body><h1>Cubeflow in your area!</h1></body></html>"
		c.String(http.StatusOK, htmlResponse)
	}
}
