package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if req.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Hello There...",
		}, nil
	}

	data, err := GetUpdates(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	logUpdate(data)

	// Handle start command
	if strings.Contains(data.Message.Text, "start") {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
		}, nil
	}

	// Handle Image
	if data.Message.Photo != nil {

		for _, image := range data.Message.Photo {
			imagePath, err := GetImagePath(image.FileID)

			if err != nil {
				log.Println(err)
			} else {
				DetectTextFromImage(imagePath)
			}
		}

		log.Println("Image detected!")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Image detected! " + data.Message.Photo[0].FileID,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil

}

func logUpdate(data Data) {
	jsonString, err := json.Marshal(data)
	if err == nil {
		log.Println(string(jsonString))
	}
}

func main() {
	lambda.Start(handler)
}
