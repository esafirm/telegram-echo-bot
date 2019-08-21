package main

import (
	"encoding/json"
	"log"
	"os"
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

	PrintJson("Updates: ", data)

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

func PrintJson(prefix string, data interface{}) {
	jsonString, err := json.Marshal(data)
	if err == nil {
		log.Println(prefix + string(jsonString))
	}
}

func main() {
	if len(os.Args) > 1 {
		fileID := os.Args[1]
		path, err := GetImagePath(fileID)

		if err != nil {
			log.Println(err)
			return
		}

		DetectTextFromImage(path)
	}

	lambda.Start(handler)
}
