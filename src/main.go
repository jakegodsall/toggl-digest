package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"jakegodsall.com/toggl-project/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
