package kubernetes

import (
	"log"
	"os"
	"path/filepath"

	env "cubeflow/pkg/config"

	argoclientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetKubeConfigClient() (*argoclientset.Clientset, *kubernetes.Clientset, error) {
	config := rest.Config{}
	cfg := &config
	var err error

	if env.Variable.Environment == "" {
		log.Fatalf("Please define ENV variable in config")

	} else if env.Variable.Environment == "development" {
		var kubeconfig string

		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			kubeconfig = os.Getenv("KUBECONFIG")
		}

		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}

	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to create in cluster config: %v", err)
			return nil, nil, err
		}
	}

	argoClientSet, err := argoclientset.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to create Argo client: %v", err)
		return nil, nil, err
	}

	kubeClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
		return nil, nil, err
	}

	return argoClientSet, kubeClientSet, nil
}
