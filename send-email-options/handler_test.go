package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type testCase struct {
	Input  string
	Output string
}

func TestHandler(t *testing.T) {
	testCases := []testCase{
		testCase{Input: "https://laurensettembrino.com", Output: "https://laurensettembrino.com"},
		testCase{Input: "https://www.laurensettembrino.com", Output: "https://www.laurensettembrino.com"},
		testCase{Input: "something random", Output: ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Origin is %s", tc.Input), func(t *testing.T) {
			headers := make(map[string]string)
			headers["Origin"] = tc.Input

			h := handler{}
			res, _ := h.Run(events.APIGatewayProxyRequest{Headers: headers})

			if res.StatusCode != 204 {
				t.Fatalf("status code should be 204, got %v", res.StatusCode)
			} else if res.Headers["Access-Control-Allow-Origin"] != tc.Output {
				t.Fatalf("expected %v for Access-Control-Allow-Origin, got %v", tc.Input, tc.Output)
			}
		})
	}
}
