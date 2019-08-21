package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

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

type GetFileResult struct {
	Result struct {
		FileID   string `json:"file_id"`
		FileSize int    `json:"file_size"`
		FilePath string `json:"file_path"`
	} `json:"result"`
}

var TELEGRAM_BOT_TOKEN string 

var client = &http.Client{}

func SendResponse(text string, chatID string) (events.APIGatewayProxyResponse, error) {
	responseData := map[string]interface{}{
		"text":    text,
		"chat_id": chatID,
	}

	var err error
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

func GetImagePath(fileID string) (string, error) {

	var filePath string

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", TELEGRAM_BOT_TOKEN, fileID)

	log.Println("Endpoint:" + endpoint)

	req, err := http.NewRequest("POST", endpoint, nil)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return filePath, err
	}

	result := GetFileResult{}
	err = getJson(req, &result)

	PrintJson("FileResut: ", result)

	filePath = fmt.Sprintf("https://api.telegram.org/bot%s/%s", TELEGRAM_BOT_TOKEN, result.Result.FilePath)

	log.Println("FilePath: " + filePath)

	return filePath, nil
}

func GetUpdates(req events.APIGatewayProxyRequest) (Data, error) {
	data := Data{}

	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		log.Println(err)
		return data, err
	}

	return data, nil
}

func getJson(req *http.Request, target interface{}) error {
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
