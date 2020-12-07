package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	sess := session.Must(session.NewSession())
	h := handler{ssmSvc: ssm.New(sess)}
	lambda.Start(h.Run)
}
