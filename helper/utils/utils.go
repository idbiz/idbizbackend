package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GithubUploadRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
}

// GithubConfig represents the GitHub credentials stored in the "ght" collection
type GithubConfig struct {
	GithubToken string `bson:"github_token"`
	GithubOwner string `bson:"github_owner"`
	GithubRepo  string `bson:"github_repo"`
}

// UploadToGithub uploads a file to a GitHub repository and returns the public file URL
func UploadToGithub(fileName, content string, db *mongo.Database) (string, error) {
	// Retrieve GitHub credentials from MongoDB
	var githubConfig GithubConfig
	collection := db.Collection("ght")

	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch the first document from "ght" collection
	err := collection.FindOne(ctx, bson.M{}).Decode(&githubConfig)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve GitHub credentials from database: %v", err)
	}

	// Validate credentials
	if githubConfig.GithubOwner == "" || githubConfig.GithubRepo == "" || githubConfig.GithubToken == "" {
		return "", fmt.Errorf("GitHub credentials missing in database")
	}

	// Construct GitHub API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubConfig.GithubOwner, githubConfig.GithubRepo, fileName)

	// Create request payload
	uploadRequest := struct {
		Message string `json:"message"`
		Content string `json:"content"`
	}{
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
	req.Header.Set("Authorization", "Bearer "+githubConfig.GithubToken)
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
	fileURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", githubConfig.GithubOwner, githubConfig.GithubRepo, fileName)

	return fileURL, nil
}
