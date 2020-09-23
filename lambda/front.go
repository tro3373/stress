package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tro3373/stress/cmd"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

var headers map[string]string = map[string]string{
	"Access-Control-Allow-Origin": "*"}

func ReturnOKResponse(statusCode int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       body,
		Headers:    headers,
		StatusCode: statusCode}, nil
}

func ReturnErrorResponse(statusCode int, errorMsg string, err error) (events.APIGatewayProxyResponse, error) {
	if errorMsg == "" {
		errorMsg = "An error occured."
		if err != nil {
			errorMsg = errorMsg + " : " + err.Error()
		}
	}

	var response = map[string]string{
		"message": errorMsg}

	body, jerr := json.Marshal(response)
	if jerr != nil {
		fmt.Println(jerr)
		errorMsg = errorMsg + " : " + jerr.Error()
	}

	var errorBuf bytes.Buffer
	json.HTMLEscape(&errorBuf, body)

	return events.APIGatewayProxyResponse{
		Body:       errorBuf.String(),
		Headers:    headers,
		StatusCode: statusCode}, err
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	fmt.Println(">>>>>> os.Args", os.Args)
	os.Args = []string{"front"}
	cmd.StartFrontendStressTest()
	response := map[string]string{
		"message": "hello"}

	body, jerr := json.Marshal(response)
	if jerr != nil {
		fmt.Println(jerr)
		return ReturnErrorResponse(500, "", jerr)
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)
	return ReturnOKResponse(200, buf.String())
}

func main() {
	lambda.Start(Handler)
}
