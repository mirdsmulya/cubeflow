package helpers

import "strings"

func ParseEnvLists(envLists string) []string {
	// Split the string by line breaks and remove leading/trailing spaces
	lines := strings.Split(strings.TrimSpace(envLists), "\n")

	// Trim spaces from each line and create a slice
	var envSingle []string
	for _, line := range lines {
		trimmedLine := strings.TrimPrefix(strings.TrimSpace(line), "- ")
		envSingle = append(envSingle, trimmedLine)
	}

	return envSingle
}
