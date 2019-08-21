package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
)

func DetectTextFromImage(imagePath string) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Detecting text from images…")

	ctx := context.Background()
	json := []byte(os.Getenv("CREDENTIALS"))

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON(json))
	if err != nil {
		log.Println(err)
		return err
	}

	// Download image
	response, e := http.Get(imagePath)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	// Create reader from HTTP respose
	images, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	reader := bytes.NewReader(images)

	image, err := vision.NewImageFromReader(reader)
	if err != nil {
		log.Println(err)
		return err
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Println(err)
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
