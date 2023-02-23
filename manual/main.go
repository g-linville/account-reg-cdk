package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

type Properties struct {
	ExternalId  string `json:"ExternalId,omitempty"`
	Arn         string `json:"Arn,omitempty"`
	Token       string `json:"Token,omitempty"`
	CallbackURL string `json:"CallbackURL,omitempty"`
	StackId     string `json:"StackId,omitempty"`
}

type Event struct {
	ResourceProperties Properties `json:"ResourceProperties,omitempty"`
}

type Payload struct {
	Token      string `json:"token"`
	ExternalID string `json:"externalID"`
	ARN        string `json:"arn"`
	StackID    string `json:"stackID,omitempty"`
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	input, _ := json.Marshal(event)
	fmt.Println("Input ", string(input))
	payload, err := json.Marshal(map[string]any{
		"spec": &Payload{
			Token:      event.ResourceProperties.Token,
			ExternalID: event.ResourceProperties.ExternalId,
			ARN:        event.ResourceProperties.Arn,
			StackID:    event.ResourceProperties.StackId,
		},
	})
	if err != nil {
		return "", err
	}
	resp, err := http.Post(event.ResourceProperties.CallbackURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	return string(payload), resp.Body.Close()
}

func main() {
	lambda.Start(HandleRequest)
}
