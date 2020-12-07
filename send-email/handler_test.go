package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var reqBodyVals map[string]string

type mockedSsm struct {
	ssmiface.SSMAPI
	GetParameterOutput *ssm.GetParameterOutput
	GetParameterError  error
}

func (m mockedSsm) GetParameter(in *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return m.GetParameterOutput, m.GetParameterError
}

func setup(t *testing.T) func(t *testing.T) {
	os.Setenv("FROM_ADDRESS", "dakota@test.com")
	os.Setenv("TO_ADDRESS", "dakota@test.com")
	os.Setenv("SMTP_HOST", "smtp.mail.com")
	os.Setenv("SMTP_PORT", "587")

	reqBodyVals = make(map[string]string)
	reqBodyVals["name"] = "Jane"
	reqBodyVals["email"] = "jane@test.com"
	reqBodyVals["subject"] = "Test subject"
	reqBodyVals["message"] = "Test message"

	return func(t *testing.T) {
		os.Unsetenv("FROM_ADDRESS")
		os.Unsetenv("TO_ADDRESS")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")

		reqBodyVals = make(map[string]string)
	}
}

func TestHandler(t *testing.T) {
	t.Run("No from address", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("FROM_ADDRESS")

		h := handler{ssmSvc: mockedSsm{}}
		res, _ := h.Run(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing from address" {
			t.Fatalf("response body should be \"Missing from address\", got %v", res.Body)
		}
	})

	t.Run("No to address", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("TO_ADDRESS")

		h := handler{ssmSvc: mockedSsm{}}
		res, _ := h.Run(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing to address" {
			t.Fatalf("response body should be \"Missing to address\", got \"%v\"", res.Body)
		}
	})

	t.Run("No smtp host", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("SMTP_HOST")

		h := handler{ssmSvc: mockedSsm{}}
		res, _ := h.Run(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing host" {
			t.Fatalf("response body should be \"Missing host\", got \"%v\"", res.Body)
		}
	})

	t.Run("No smtp port", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("SMTP_PORT")

		h := handler{ssmSvc: mockedSsm{}}
		res, _ := h.Run(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing port" {
			t.Fatalf("response body should be \"Missing port\", got \"%v\"", res.Body)
		}
	})

	t.Run("Request body is not valid JSON", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		h := handler{ssmSvc: mockedSsm{}}
		res, _ := h.Run(events.APIGatewayProxyRequest{Body: "{\"test: wut}"})

		if res.StatusCode != 400 {
			t.Fatalf("status code should be 400, got %v", res.StatusCode)
		} else if res.Body != "Unable to unmarshal request body" {
			t.Fatalf("response body should be \"Unable to unmarshal request body\", got \"%v\"", res.Body)
		}
	})

	for _, key := range []string{"name", "email", "subject", "message"} {
		t.Run(fmt.Sprintf("Request body is missing %v field", key), func(t *testing.T) {
			teardown := setup(t)
			defer teardown(t)
			delete(reqBodyVals, key)
			reqBody, err := json.Marshal(reqBodyVals)
			if err != nil {
				t.Fatal("Unable to marshal request body")
			}

			h := handler{ssmSvc: mockedSsm{}}
			res, _ := h.Run(events.APIGatewayProxyRequest{Body: string(reqBody)})

			expectedResBody := fmt.Sprintf("Missing %s", key)
			if res.StatusCode != 400 {
				t.Fatalf("status code should be 400, got %v", res.StatusCode)
			} else if res.Body != expectedResBody {
				t.Fatalf("response body should be \"%v\", got \"%v\"", expectedResBody, res.Body)
			}
		})
	}

	t.Run("SSM returns an error", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)
		reqBody, err := json.Marshal(reqBodyVals)
		if err != nil {
			t.Fatal("Unable to marshal request body")
		}

		mSsm := mockedSsm{GetParameterError: errors.New("Oops")}
		h := handler{ssmSvc: mSsm}
		res, _ := h.Run(events.APIGatewayProxyRequest{Body: string(reqBody)})

		expectedResBody := "Unable to get email password"
		if res.StatusCode != 502 {
			t.Fatalf("status code should be 502, got %v", res.StatusCode)
		} else if res.Body != expectedResBody {
			t.Fatalf("response body should be \"%v\", got \"%v\"", expectedResBody, res.Body)
		}
	})
}
