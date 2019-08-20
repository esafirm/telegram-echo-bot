package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var TELEGRAM_BOT_TOKEN string

type PhotoSize struct {
	FileID string `json:"file_id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Data struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Photo []PhotoSize `json:"photo"`
		Date  int         `json:"date"`
		Text  string      `json:"text"`
	} `json:"message"`
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if req.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Hello There...",
		}, nil
	}

	var err error

	data := Data{}

	if err = json.Unmarshal([]byte(req.Body), &data); err != nil {

		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	}

	log.Println(data)

	if strings.Contains(data.Message.Text, "start") {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
		}, nil
	}

	if data.Message.Photo != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Image detected! " + data.Message.Photo[0].FileID,
		}, nil
	}

	responseData := map[string]interface{}{
		"text":    data.Message.Text,
		"chat_id": data.Message.Chat.ID,
	}

	var responseDataJSON []byte

	if responseDataJSON, err = json.Marshal(responseData); err != nil {

		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	}

	if request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_BOT_TOKEN),
		bytes.NewReader(responseDataJSON)); err != nil {

		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	} else {

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		if _, err = client.Do(request); err != nil {
			log.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: 503,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil

}

func main() {
	// detectText()
	lambda.Start(handler)
}

func detectText() error {
	file := "./gojek.png"

	ctx := context.Background()
	json := []byte(os.Getenv("CREDENTIALS"))

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON(json))
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Println("No text found.")
	} else {
		fmt.Println("Text:")
		for _, annotation := range annotations {
			fmt.Println("%q\n", annotation.Description)
		}
	}

	return nil
}

func checkVision() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.
	filename := "./gojek.png"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Println(label.Description)
	}
}
