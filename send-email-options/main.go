package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	h := handler{}
	lambda.Start(h.Run)
}
