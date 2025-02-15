package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type TogglClient struct {
	AuthHeader string
}

func NewTogglClient(authHeader string) *TogglClient {
	return &TogglClient{AuthHeader:  authHeader}
}

func (client *TogglClient) GetProjectMap() (map[string]string, error) {
	workspaceId := os.Getenv("TOGGL_WORKSPACE_ID")
	if workspaceId == "" {
		return nil, errors.New("environment variable TOGGL_WORKSPACE_ID is not set")
	}
	
	apiURL := "https://api.track.toggl.com/api/v9/workspaces/" + workspaceId + "/projects"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", client.AuthHeader)
	resp, err := http.DefaultClient.Do(req)
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