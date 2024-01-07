package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetArgoCDToken(argoCDURL string, username string, password string) (string, error) {

	// Define the API endpoint for login
	loginEndpoint := fmt.Sprintf("%s/api/v1/session", argoCDURL)

	// Create an HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // This skips certificate verification
		},
	}
	// Create a POST request to log in
	loginData := map[string]interface{}{
		"username": username,
		"password": password,
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		return "loginJSON ", err
	}

	req, err := http.NewRequest("POST", loginEndpoint, bytes.NewReader(loginJSON))
	if err != nil {
		return "req NewRequest ", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the POST request to log in
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// Handle the error response
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to log in to ArgoCD: %s", responseBody)
	}

	// Extract the token from the response
	var loginResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return "", err
	}

	token, ok := loginResponse["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in login response")
	}

	return token, nil
}
