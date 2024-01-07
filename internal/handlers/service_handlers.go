package handlers

import (
	"log"
	"net/http"

	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"

	"cubeflow/internal/argo"
	kubeFunc "cubeflow/internal/kubernetes"
	"cubeflow/internal/slack"

	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error message: %s", err)
		}
		log.Printf("deployment %s in namespace %s\n", deployments, namespace)
		c.JSON(http.StatusOK, "done, check terminal deh")
	}
}

func RestartPodsV1(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		rolloutName := c.Param("rolloutName")

		channelName, found := slack.GetChannelNameByAppName(rolloutName)
		if !found {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channel name"})
			return
		}

		err := argo.RestartArgoRollout(clientset, argoClientset, namespace, rolloutName)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			slack.SendSlackMessage(channelName, "Oops its failed, please check again your webhook path value :eyes:")
			return
		}
		c.JSON(http.StatusOK, "Rollout restart initiated!")
		slack.SendSlackMessage(channelName, "Processing restart for "+rolloutName+"! \nPlease wait a second \nTriggered by: Alert system")
	}
}

func RestartPodsV2(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, appEnv string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, rolloutName, triggerUserName, incomingSlackChannel := slack.AuthChecker(c)
		ChannelMatched := slack.IsChannelMatched(incomingSlackChannel, rolloutName)
		if !authenticated || !ChannelMatched {
			log.Printf("Error message: failed to authenticated or wrong channel or wrong application name")
			c.JSON(http.StatusOK, "Failed, you might be running this command in wrong channel or wrong application name")
			return
		}

		// change to only trigger restart when all the apps name is valid ?
		namespace := c.Param("namespace")
		proceedList, errList := argo.MultipleRestartPods(clientset, argoClientset, namespace, rolloutName)
		if errList != "" {
			rolloutNameHelper := argo.GetRolloutListParsed(clientset, argoClientset, namespace)
			c.JSON(http.StatusOK, "Restart failed for "+errList+", please check again your prompt value (options: "+rolloutNameHelper+")")
			if proceedList != "" {
				slack.SendSlackMessage(incomingSlackChannel, proceedList+" restart on "+appEnv+" cluster: partial success :rocket:\nNamespace: "+namespace+"\nTriggered by: "+triggerUserName)
			}
			return
		}

		c.JSON(http.StatusOK, "Rollout restart initiated!")
		slack.SendSlackMessage(incomingSlackChannel, proceedList+" restart on "+appEnv+" cluster: success :rocket:\nNamespace: "+namespace+"\nTriggered by: "+triggerUserName)
	}
}

func ScaleDownServiceV1(clientset *kubernetes.Clientset, argoClientset *argoclientset.Clientset, namespaces []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "Service scale down initiated!")

		log.Println("Scale down all Rollouts Initiated")
		// Scale down all Rollouts in the namespace
		err := argo.ScaleDownAllRolloutsNamespaces(argoClientset, namespaces)
		if err != nil {
			log.Println("Error scale down Rollout: ", err)
		}

		log.Println("Scale down all Deployments Initiated")
		// Scale down all Deployments in the namespace
		err = kubeFunc.ScaleDownAllDeploymentsNamespaces(clientset, namespaces)
		if err != nil {
			log.Println("Error scale down deployment: ", err)
		}

		log.Println("Scale down all Statefulsets Initiated")
		// Scale down all Statefulsets in the namespace
		err = kubeFunc.ScaleDownAllStatefulsetsNamespaces(clientset, namespaces)
		if err != nil {
			log.Println("Error scale down Statefulset: ", err)
		}

		log.Println("Scaled down all Service completed")
	}
}

func RollbackService(argoClientset *argoclientset.Clientset, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, inputPrompt, triggerUserName, incomingSlackChannel := slack.AuthChecker(c)
		ChannelMatched := slack.IsChannelMatched(incomingSlackChannel, inputPrompt)
		if !authenticated || !ChannelMatched {
			log.Printf("Error message: failed to authenticated or wrong channel or wrong application name")
			c.JSON(http.StatusOK, "Failed, you might be running this command in wrong channel or wrong application name")
			return
		}

		namespace := c.Param("namespace")
		newTag, prevTag, serviceName, err := argo.RollbackVersion(argoClientset, namespace, inputPrompt, env)
		if err != nil || newTag == prevTag {
			c.JSON(http.StatusOK, "Rollback service failed! please check again your prompt value (e.g. service-name,v0.0.0-production)")
			slack.SendSlackMessage(incomingSlackChannel, "Rollback deployment failed :heavy_multiplication_x: \nTriggered by: "+triggerUserName)
			log.Printf("Rollback failed: %v", err)
			return
		}
		c.JSON(http.StatusOK, "Rollback service initiated!")
		slack.SendSlackMessage(incomingSlackChannel, "Rollback deployment for "+serviceName+" on "+env+" cluster: success :recycle:\nNamespace: "+namespace+"\nCurrent: "+newTag+" \nPrevious: "+prevTag+"\nTriggered by: "+triggerUserName)
	}
}
