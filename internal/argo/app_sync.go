package argo

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SyncArgoCDApplication(token string, applicationName string, argoCDURL string) error {
	apiEndpoint := fmt.Sprintf("%s/api/v1/applications/%s/sync", argoCDURL, applicationName)

	// Create an HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // This skips certificate verification
		},
	}

	// Create a POST request to trigger the synchronization
	req, err := http.NewRequest("POST", apiEndpoint, nil)
	if err != nil {
		log.Println("Error when Create a POST request ", err)
		return err
	}

	// Send the POST request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error when send the post request ", err)
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to sync ArgoCD Application: %s", responseBody)
	}

	log.Printf("ArgoCD Application '%s' synced successfully.\n", applicationName)
	return nil
}
