package argo

import (
	"context"
	helperFunc "cubeflow/internal/helpers"
	"errors"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RollbackVersion(argoClientset *argoclientset.Clientset, namespace string, inputPrompt string, env string) (string, string, string, error) {
	var newVersionImage string
	decodedInputPrompt, _ := url.QueryUnescape(inputPrompt)
	rolloutName, newTag := SplitInputPrompt(decodedInputPrompt)

	if newTag == "" {
		log.Println("Missing second expression to rollback 'rolloutName,versionTag'")
		err := errors.New("error when splitting prompt input")
		return newTag, newTag, rolloutName, err
	}

	if !checkRegexTagVersion(newTag, env) {
		log.Println("Error: tag format not meet the standards format (e.g. v0.0.1-alpha)")
		err := errors.New("error: new tag format is not standard")
		return newTag, newTag, rolloutName, err
	}

	rolloutManifest, err := GetRolloutManifest(argoClientset, namespace, rolloutName)
	if err != nil {
		log.Println("Error when get rollout manifest, service name might be wrong 'roloutName,versionTag'")
		return newTag, newTag, rolloutName, err
	}

	roloutContainerIndex := ContainerIndexCheck(rolloutManifest, rolloutName)
	prevTag := GetTag(rolloutManifest, roloutContainerIndex)
	newVersionImage = CompileNewImageTag(rolloutManifest, roloutContainerIndex, newTag)
	rolloutManifest.Spec.Template.Spec.Containers[roloutContainerIndex].Image = newVersionImage

	err = updateRollout(argoClientset, namespace, rolloutManifest)
	if err != nil {
		log.Printf("Error when applying rollout manifest: %v", err)
		return newTag, prevTag, rolloutName, err
	}

	return newTag, prevTag, rolloutName, nil
}

func SplitInputPrompt(inputPrompt string) (string, string) {
	splitPatameter := ","
	if !helperFunc.MoreThanOneExpressionCheck(inputPrompt, splitPatameter) {
		log.Println("Expected two input value separate with ',' ")
		return "", ""
	}

	inputPromptSplitted := helperFunc.ParseMoreThanOneExpression(inputPrompt, splitPatameter)
	rolloutName := inputPromptSplitted[0]
	rolbackVersionTag := inputPromptSplitted[1]
	return rolloutName, rolbackVersionTag
}

func checkRegexTagVersion(inputString string, env string) bool {
	var pattern string
	if env == "production" {
		pattern = `^v\d+\.\d+\.\d+-production$`
		re := regexp.MustCompile(pattern)
		return re.MatchString(inputString)
	}
	pattern = `^v\d+\.\d+\.\d+-alpha$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(inputString)
}

func CompileNewImageTag(rolloutManifest *v1alpha1.Rollout, index int, rolbackVersionTag string) string {
	imagePathWithVersion := rolloutManifest.Spec.Template.Spec.Containers[index].Image
	imagePathWithVersionArr := strings.Split(imagePathWithVersion, ":")
	imagePath := imagePathWithVersionArr[0]
	newRollbackVersionImage := imagePath + ":" + rolbackVersionTag
	return newRollbackVersionImage
}

func updateRollout(argoClientset *argoclientset.Clientset, namespace string, rolloutObj *v1alpha1.Rollout) error {
	_, err := argoClientset.ArgoprojV1alpha1().Rollouts(namespace).Update(context.TODO(), rolloutObj, metav1.UpdateOptions{})
	if err != nil {
		log.Println("Error when Apply the changes to the Argo Rollout")
		return err
	}
	return nil
}
