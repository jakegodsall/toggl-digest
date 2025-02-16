package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type TogglClient struct {
	AuthHeader string
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TimeEntry struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Start       string `json:"start"`
	End         string `json:"stop"`
}

type ProjectTime struct {
	ID          int
	ProjectName string
	Description string
	Duration    int
	Start       string
	End         string
}

func NewTogglClient(authHeader string) *TogglClient {
	return &TogglClient{AuthHeader:  authHeader}
}

func (client *TogglClient) GetProjectMap() (map[int]string, error) {
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

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode the JSON response: %w", err)
	}

	projectMap := make(map[int]string)
	for _, project := range projects {
		projectMap[project.ID] = project.Name
	}

	return projectMap, nil
}

func (client *TogglClient) GetTimeEntries() ([]TimeEntry, error) {
	workspaceId := os.Getenv("TOGGL_WORKSPACE_ID")
	if workspaceId == "" {
		return nil, errors.New("environment variable TOGGL_WORKSPACE_ID is not set")
	}

	apiUrl := "https://api.track.toggl.com/api/v9/me/time_entries"
	req, err := http.NewRequest("GET", apiUrl, nil)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var timeEntries []TimeEntry

	json.Unmarshal(bodyBytes, &timeEntries)
	err = json.Unmarshal(bodyBytes, &timeEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return timeEntries, nil
}

func (client *TogglClient) GetTimeEntriesWithProjects(timeEntries []TimeEntry, projectMap map[int]string) []ProjectTime {
	var projectTimes []ProjectTime

	for _, timeEntry := range timeEntries {
		// Check if the projectId exists in the map
		if projectName, exists := projectMap[timeEntry.ProjectID]; exists {
			projectTimes = append(projectTimes, ProjectTime{
				ID:          timeEntry.ID,
				ProjectName: projectName,
				Description: timeEntry.Description,
				Duration:    timeEntry.Duration,
				Start:       timeEntry.Start,
				End:         timeEntry.End,
			})
		}
	}

	return projectTimes
}