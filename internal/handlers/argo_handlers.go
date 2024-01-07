package handlers

import (
	"cubeflow/internal/argo"
	"cubeflow/internal/helpers"
	kubeFunc "cubeflow/internal/kubernetes"
	"cubeflow/internal/slack"
	"cubeflow/pkg/config"
	"log"
	"net/http"
	"time"

	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

func SyncArgoApp(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, appEnv string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, applicationName, triggerUserName, incomingSlackChannel := slack.AuthChecker(c)
		ChannelMatched := slack.IsChannelMatched(incomingSlackChannel, applicationName)
		if !authenticated || !ChannelMatched {
			log.Printf("Error message: failed to authenticated or wrong channel or wrong application name")
			c.JSON(http.StatusOK, "Failed, you might be running this command in wrong channel or wrong application name")
			return
		}

		token, err := argo.GetArgoCDToken(config.Variable.ArgoCD.ServerURL, config.Variable.ArgoCD.Username, config.Variable.ArgoCD.Password)
		if err != nil {
			log.Println("Failed to obtain authentication token: ", err)
			c.JSON(http.StatusOK, "Error, make sure your argocd credential config is correct")
			return
		}

		go func() {
			err = argo.SyncArgoCDApplication(token, applicationName, config.Variable.ArgoCD.ServerURL)
			if err != nil {
				slack.SendSlackMessage(incomingSlackChannel, "Oops sync failed, please check again your prompt value :eyes: \nTriggered by: "+triggerUserName)
				log.Println("Failed to sync ArgoCD Application: ", err)
				return
			}
			// Delay notification for 20 seconds
			duration := 20 * time.Second
			time.Sleep(duration)
			slack.SendSlackMessage(incomingSlackChannel, "Syncing latest "+applicationName+" tag on "+appEnv+": success :white_check_mark: \nPlease wait a second.\nTriggered by: "+triggerUserName)
		}()
		c.JSON(http.StatusOK, "Sync initiated! please wait a second")
	}
}

func PromoteRollout(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, appEnv string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, rolloutName, triggerUserName, incomingSlackChannel := slack.AuthChecker(c)
		ChannelMatched := slack.IsChannelMatched(incomingSlackChannel, rolloutName)
		if !authenticated || !ChannelMatched {
			log.Printf("Error message: failed to authenticated or wrong channel or wrong application name")
			c.JSON(http.StatusOK, "Failed, you might be running this command in wrong channel or wrong application name")
			return
		}

		namespace := c.Param("namespace")
		if rolloutName == "status" {
			rolloutList, _ := argo.GetRolloutList(clientset, argoClientset, namespace)
			rolloutStatusList := argo.GetRolloutStatusList(rolloutList)
			rolloutStatusListStr := argo.ParseSingleLineArr(rolloutStatusList)
			c.JSON(http.StatusOK, "Pre-prod readiness -- "+rolloutStatusListStr)
			return
		}

		rolloutManifest, err := argo.GetRolloutManifest(argoClientset, namespace, rolloutName)
		if err != nil {
			rolloutList, _ := argo.GetRolloutList(clientset, argoClientset, namespace)
			allRolloutname := argo.GetAllRolloutName(rolloutList)
			strRolloutName := helpers.ParseToSingleLine(allRolloutname, ", ")
			c.JSON(http.StatusOK, "Promote failed, please check again your prompt value (options: "+strRolloutName+")")
			return
		}

		podsCrashed, _, err := kubeFunc.IsPodsCrashed(clientset, namespace, rolloutName)
		if err != nil {
			c.JSON(http.StatusOK, "Cubeflow can't check status of "+rolloutName+", please check cubeflow logs for details")
			return
		}

		if podsCrashed {
			c.JSON(http.StatusOK, "One or more pods of "+rolloutName+" is crashed/not ready, please check again before deployment")
			return
		}

		isReadyToPromote, promotionStatus := argo.GetPromoteStatus(rolloutManifest)
		if !isReadyToPromote {
			c.JSON(http.StatusOK, promotionStatus+" ps: check '/{promote command} status' ")
			return
		}

		err = argo.PromoteRollout(argoClientset, rolloutManifest, namespace)
		if err != nil {
			c.JSON(http.StatusOK, "Promote failed, system can't update rollout object")
			return
		}

		log.Printf("Promoting Argo Rollout %s in namespace %s\n", rolloutName, namespace)
		c.JSON(http.StatusOK, "Promote initiated! please wait 30s until it deploys")

		go func() {
			// Delay notification for 30 seconds
			duration := 32 * time.Second
			time.Sleep(duration)
			rolloutManifest, _ := argo.GetRolloutManifest(argoClientset, namespace, rolloutName)
			readyToPromote, _ := argo.GetPromoteStatus(rolloutManifest)
			if readyToPromote {
				slack.SendSlackMessage(incomingSlackChannel, "Promoting "+rolloutName+" to "+appEnv+" was failed! :x: \nTriggered by: "+triggerUserName)
				return
			}
			latestImage := argo.GetImage(rolloutManifest)
			slack.SendSlackMessage(incomingSlackChannel, "Promoting "+rolloutName+" to "+appEnv+": success! :white_check_mark: \nTag version: "+latestImage+"\nTriggered by: "+triggerUserName)
		}()
	}
}
