package main

import (
	"context"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
)

func DetectTextFromImage(imagePath string) error {


	ctx := context.Background()
	json := []byte(os.Getenv("CREDENTIALS"))

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON(json))
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(imagePath)
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		log.Println("No text found.")
	} else {
		log.Println("Text:")
		for _, annotation := range annotations {
			log.Println("%q\n", annotation.Description)
		}
	}

	return nil
}
