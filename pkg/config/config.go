package config

import (
	"os"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"env"`
	ArgoCD      struct {
		ServerURL       string   `yaml:"server_url"`
		Username        string   `yaml:"username"`
		Password        string   `yaml:"password"`
		ApplicationName []string `yaml:"appname"`
	} `yaml:"argocd"`
	Database struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`
	Slack struct {
		SigningSecret string          `yaml:"signing_secret"`
		Token         string          `yaml:"token"`
		Channels      []SlackChannels `yaml:"channel"`
	} `yaml:"slack"`
	GCP struct {
		ProjectID       string `yaml:"project_id"`
		CredentialsPath string `yaml:"credentials_path"`
		BucketName      string `yaml:"bucket_name"`
		GKE             struct {
			ClusterName      string   `yaml:"cluster_name"`
			Zone             string   `yaml:"zone"`
			NamespaceToScale []string `yaml:"namespace_to_scale"`
		} `yaml:"gke"`
	} `yaml:"gcp"`
}

type SlackChannels struct {
	Name        string   `yaml:"name"`
	RolloutName []string `yaml:"rollouts_name"`
	ArgoAppName []string `yaml:"argo_app_name"`
}

var Variable *Config

func Load(src string) error {
	if Variable != nil {
		log.Info("Config successfully loaded")
		return nil
	}
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&Variable); err != nil {
		return err
	}
	return nil
}
