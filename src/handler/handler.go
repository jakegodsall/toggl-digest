package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"jakegodsall.com/toggl-project/auth"
	"jakegodsall.com/toggl-project/toggl"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader, err := auth.GetAuthHeaderValue()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error getting auth header: %s", err),
			StatusCode: 500,
		}, nil
	}

	client := toggl.NewTogglClient(authHeader)

	projectMap, err := client.GetProjectMap()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error getting project map.",
			StatusCode: 500,
		}, nil
	}
	fmt.Println(projectMap)

	client.GetTimeEntries()

	// Fetch user info
	apiURL := "https://api.track.toggl.com/api/v9/me"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error creating request: %s", err),
			StatusCode: 500,
		}, nil
	}

	req.Header.Add("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error making request: %s", err),
			StatusCode: 500,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error reading response: %s", err),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}
