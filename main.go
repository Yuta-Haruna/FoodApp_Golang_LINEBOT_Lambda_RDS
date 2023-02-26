package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func hello(event MyEvent) (MyResponse, error) {
	return MyResponse{Message: fmt.Sprintf("Hello %s!!", event.Name)}, nil
}

func main() {
	lambda.Start(hello)
}
