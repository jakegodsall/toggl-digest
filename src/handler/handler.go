package handler

import (
	"fmt"

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

	timeEntries, err := client.GetTimeEntries()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error getting time entries.",
			StatusCode: 500,
		}, nil
	}

	client.GetTimeEntriesWithProjects(timeEntries, projectMap)

	return events.APIGatewayProxyResponse{
		Body:       "success",
		StatusCode:  200,
	}, nil
}
