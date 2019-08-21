package main

import (
	"context"
	"fmt"
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
		fmt.Println("No text found.")
	} else {
		fmt.Println("Text:")
		for _, annotation := range annotations {
			fmt.Println("%q\n", annotation.Description)
		}
	}

	return nil
}
