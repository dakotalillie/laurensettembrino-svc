package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// EmailConfig is a struct which contains the requisite values for sending an email.
type EmailConfig struct {
	From        string // The email address the email will be sent from
	To          string // The email address the email will be sent to
	Host        string // The SMTP host for the email server
	Port        string // The SMTP port for the email server
	Password    string // The password for the sending account
	SenderName  string // The name of the person who is sending the email
	SenderEmail string // The email of the person who is sending the email, used for the reply to
	Subject     string // The subject of the email
	Message     string // The message body of the email
}

// NewEmailConfig creates a new email config, drawing values from environment variables, the request
// body, and SSM in AWS.
func NewEmailConfig(reqBody string, ssmSvc ssmiface.SSMAPI) (EmailConfig, error) {
	config := EmailConfig{}

	from := os.Getenv("FROM_ADDRESS")
	if from == "" {
		return config, &codedError{Code: 500, Message: "Missing from address"}
	}

	to := os.Getenv("TO_ADDRESS")
	if to == "" {
		return config, &codedError{Code: 500, Message: "Missing to address"}
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return config, &codedError{Code: 500, Message: "Missing host"}
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		return config, &codedError{Code: 500, Message: "Missing port"}
	}

	parsedBody := make(map[string]string)
	err := json.Unmarshal([]byte(reqBody), &parsedBody)
	if err != nil {
		return config, &codedError{Code: 400, Message: "Unable to unmarshal request body"}
	}

	for _, key := range []string{"name", "email", "subject", "message"} {
		if parsedBody[key] == "" {
			return config, &codedError{Code: 400, Message: "Missing " + key}
		}
	}

	res, err := ssmSvc.GetParameter(&ssm.GetParameterInput{Name: aws.String("/LaurenSettembrino/EMAIL_PASSWORD"), WithDecryption: aws.Bool(true)})
	if err != nil {
		return config, &codedError{Code: 502, Message: "Unable to get email password"}
	}
	password := *res.Parameter.Value

	config.From = from
	config.To = to
	config.Host = host
	config.Port = port
	config.Password = password
	config.SenderName = parsedBody["name"]
	config.SenderEmail = parsedBody["email"]
	config.Subject = parsedBody["subject"]
	config.Message = parsedBody["message"]

	return config, nil
}
