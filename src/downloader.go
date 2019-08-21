package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadImage(url string) (*os.File, error){
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// don't worry about errors
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("temp.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")

	return file, err
}
