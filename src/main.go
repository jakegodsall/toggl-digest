package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := os.Getenv("TOGGL_EMAIL")
	password := os.Getenv("TOGGL_PASSWORD")

	auth := email + ":" + password

	fmt.Println(auth)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	authHeader := "Basic " + authEncoded

	
	// Define the GET request
	apiURL := "https://api.track.toggl.com/api/v9/me"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error creating request: %s", err.Error()),
			StatusCode: 500,
			}, nil
		}
		
		req.Header.Add("Authorization", authHeader)
		
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
