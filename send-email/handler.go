package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
)

type handler struct{}

func (h *handler) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// These values were copied from here: https://docs.aws.amazon.com/apigateway/latest/developerguide/how-to-cors.html
	headers := make(map[string]string)
	headers["Access-Control-Allow-Headers"] = "Content-Type"
	headers["Access-Control-Allow-Methods"] = "OPTIONS,POST"
	headers["Access-Control-Max-Age"] = "300"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Vary"] = "Origin"

	requestOrigin := request.Headers["origin"]
	if requestOrigin == "" {
		requestOrigin = request.Headers["Origin"]
	}

	allowedOrigins := []string{"https://laurensettembrino.com", "https://www.laurensettembrino.com"}
	for _, origin := range allowedOrigins {
		if requestOrigin == origin {
			headers["Access-Control-Allow-Origin"] = origin
			break
		}
	}

	if headers["Access-Control-Allow-Origin"] == "" {
		return events.APIGatewayProxyResponse{StatusCode: 403, Headers: headers, Body: "Invalid origin"}, nil
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Headers: headers, Body: err.Error()}, nil
	}

	config, err := NewEmailConfig(ctx, request.Body, cfg)
	if codedError, ok := err.(*codedError); ok {
		return events.APIGatewayProxyResponse{StatusCode: codedError.Code, Headers: headers, Body: codedError.Error()}, nil
	}

	mailer := Mailer{config}
	if err = mailer.Send(); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 502, Headers: headers, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers}, nil
}
