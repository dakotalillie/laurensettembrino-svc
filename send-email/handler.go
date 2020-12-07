package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type handler struct {
	ssmSvc ssmiface.SSMAPI
}

func (h *handler) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// These values were copied from here: https://docs.aws.amazon.com/apigateway/latest/developerguide/how-to-cors.html
	headers := make(map[string]string)
	headers["Access-Control-Allow-Headers"] = "Content-Type"
	headers["Access-Control-Allow-Methods"] = "OPTIONS,POST"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	allowedOrigins := []string{"https://laurensettembrino.com", "https://www.laurensettembrino.com"}
	for _, origin := range allowedOrigins {
		if request.Headers["Origin"] == origin {
			headers["Access-Control-Allow-Origin"] = origin
			break
		}
	}

	config, err := NewEmailConfig(request.Body, h.ssmSvc)
	if codedError, ok := err.(*codedError); ok {
		return events.APIGatewayProxyResponse{StatusCode: codedError.Code, Headers: headers, Body: codedError.Error()}, nil
	}

	mailer := Mailer{config}
	if err = mailer.Send(); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 502, Headers: headers, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers}, nil
}
