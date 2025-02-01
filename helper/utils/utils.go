package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GithubUploadRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
}

// UploadToGithub uploads a file to a GitHub repository and returns the public file URL
func UploadToGithub(fileName, content string) (string, error) {
	// Load environment variables from .env file
	err := godotenv.Load("D:/CAMPUS/KULIAH/Smest. 5/Sistem Informasi Geografis/repos/idbizbackend/.env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	// Get owner, repo, and token from environment variables
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubToken := os.Getenv("GITHUB_TOKEN")

	// Validate environment variables
	if githubOwner == "" || githubRepo == "" || githubToken == "" {
		return "", fmt.Errorf("GITHUB_OWNER, GITHUB_REPO, or GITHUB_TOKEN is not set in the environment")
	}

	// Construct GitHub API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubOwner, githubRepo, fileName)

	// Create request payload
	uploadRequest := GithubUploadRequest{
		Message: "Upload file " + fileName,
		Content: content,
	}

	// Convert request payload to JSON
	jsonData, err := json.Marshal(uploadRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal upload request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := ioutil.ReadAll(resp.Body)

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to upload file to GitHub: %s, response: %s", resp.Status, string(body))
	}

	// Generate public file URL
	fileURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", githubOwner, githubRepo, fileName)

	return fileURL, nil
}
