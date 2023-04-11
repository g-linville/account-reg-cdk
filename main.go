package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	RequestType        string     `json:"RequestType"`
	ResponseURL        string     `json:"ResponseURL"`
	StackId            string     `json:"StackId"`
	RequestId          string     `json:"RequestId"`
	ResourceType       string     `json:"ResourceType"`
	LogicalResourceId  string     `json:"LogicalResourceId"`
	ResourceProperties Properties `json:"ResourceProperties"`
}

type Payload struct {
	Token      string `json:"token"`
	ExternalID string `json:"externalID"`
	ARN        string `json:"arn"`
	StackID    string `json:"stackID,omitempty"`
}

type Response struct {
	Status             string            `json:"Status"`
	PhysicalResourceId string            `json:"PhysicalResourceId"`
	StackId            string            `json:"StackId"`
	RequestId          string            `json:"RequestId"`
	LogicalResourceId  string            `json:"LogicalResourceId"`
	Data               map[string]string `json:"Data"`
}

type Metadata struct {
	UID  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
}

type APIResponse struct {
	Metadata `json:"metadata,omitempty"`
}

func logEvent(event Event) {
	if len(event.ResourceProperties.Token) > 5 {
		event.ResourceProperties.Token = event.ResourceProperties.Token[:5]
	}
	input, _ := json.Marshal(event)
	fmt.Println("Input ", string(input))
}

func handleResponse(event Event, r Response) (string, error) {
	eventResponse, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	fmt.Printf("Uploading to [%s] payload: %s\n", event.ResponseURL, string(eventResponse))
	req, err := http.NewRequest(http.MethodPut, event.ResponseURL, bytes.NewBuffer(eventResponse))
	if err != nil {
		return "", err
	}
	req.ContentLength = int64(len(eventResponse))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return string(eventResponse), err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload response got [%d]: %s", resp.StatusCode, string(data))
	}

	return string(eventResponse), err
}

func HandleRequest(_ context.Context, event Event) (string, error) {
	logEvent(event)

	if event.RequestType != "Create" {
		return handleResponse(event, Response{
			Status:             "SUCCESS",
			StackId:            event.StackId,
			RequestId:          event.RequestId,
			PhysicalResourceId: "unknown",
			LogicalResourceId:  event.LogicalResourceId,
		})
	}
	if event.ResourceProperties.Token == "" {
		return handleResponse(event, Response{
			Status:             "SUCCESS",
			StackId:            event.StackId,
			RequestId:          event.RequestId,
			PhysicalResourceId: "no-token",
			LogicalResourceId:  event.LogicalResourceId,
		})
	}

	payload, err := json.Marshal(map[string]any{
		"kind":       "AWSAccountRegistrationCallback",
		"apiVersion": "internal.account.hub.acorn.io/v1",
		"metadata": map[string]any{
			"generateName": "create",
		},
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
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create resource got [%d]: %s", resp.StatusCode, string(data))
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(data, &apiResponse); err != nil {
		return "", err
	}

	return handleResponse(event, Response{
		Status:             "SUCCESS",
		PhysicalResourceId: apiResponse.UID,
		StackId:            event.StackId,
		RequestId:          event.RequestId,
		LogicalResourceId:  event.LogicalResourceId,
		Data: map[string]string{
			"Name": apiResponse.Name,
			"UID":  apiResponse.UID,
		},
	})
}

func main() {
	lambda.Start(HandleRequest)
}
