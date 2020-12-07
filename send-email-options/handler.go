package main

import (
	"github.com/aws/aws-lambda-go/events"
)

type handler struct{}

func (h *handler) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Headers"] = "Content-Type"
	headers["Access-Control-Allow-Methods"] = "OPTIONS,POST"
	headers["Access-Control-Max-Age"] = "300"
	headers["Vary"] = "Origin"

	allowedOrigins := []string{"https://laurensettembrino.com", "https://www.laurensettembrino.com"}
	for _, origin := range allowedOrigins {
		if request.Headers["Origin"] == origin {
			headers["Access-Control-Allow-Origin"] = origin
			break
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 204, Headers: headers}, nil
}
