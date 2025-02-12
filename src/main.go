package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getAuthHeaderValue() (string, error) {
	// Get the environment variables
	email := os.Getenv("TOGGL_EMAIL")
	password := os.Getenv("TOGGL_PASSWORD")
	if email == "" {
		return "", errors.New("environment variable TOGGL_EMAIL is not set")
	}
	if password == "" {
		return "", errors.New("environment variable TOGGL_PASSWORD is not set")
	}

	// Concatenate and encode to base64
	auth := email + ":" + password
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	headerValue := "Basic " + authEncoded

	return headerValue, nil
}

func getProjectMap() (map[string]string, error) {
	workspaceId := os.Getenv("TOGGL_WORKSPACE_ID")
	if workspaceId == "" {
		return nil, errors.New("environment variable TOGGL_WORKSPACE_ID is not set")
	}
	
	apiURL := "https://api.track.toggl.com/api/v9/workspaces/" + workspaceId + "/projects"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	authHeaderValue, err := getAuthHeaderValue()
	if err != nil {
		return nil, errors.New("could not create authentication header value")
	}

	req.Header.Set("Authorization", authHeaderValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var projects []struct {
		ID int
		Name string
	}
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode the JSON response: %w", err)
	}

	projectMap := make(map[string]string)
	for _, project := range projects {
		projectMap[fmt.Sprint(project.ID)] = project.Name
	}

	return projectMap, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projectMap, err := getProjectMap()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error getting project map.",
			StatusCode: 500,
		}, nil
	}
	fmt.Println(projectMap)

	authHeaderValue := getAuthHeaderValue()
	
	// Define the GET request
	apiURL := "https://api.track.toggl.com/api/v9/me"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error creating request: %s", err.Error()),
			StatusCode: 500,
			}, nil
		}
		
		req.Header.Add("Authorization", authHeaderValue)
		
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error making request: %s", err.Error()),
			StatusCode: 500,
		}, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return events.APIGatewayProxyResponse{
			Body:       "Unauthorised: Invalid Toggl credentials",
			StatusCode: 401,
		}, nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error reading response: %s", err.Error()),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
