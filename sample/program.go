package main

import (
	"log"
	"os"

	abbyysdk "github.com/microo8/cloudsdk-client-go"
)

const (
	HOST           = "https://cloud-eu.ocrsdk.com"
	APPLICATION_ID = "PASTE_APPLICATION_ID"
	PASSWORD       = "PASTE_APPLICATION_PASSWORD"

	FILEPATH = "PASTE_IMAGE_FILEPATH"
)

func main() {
	ocrClient := abbyysdk.NewOcrClient(HOST, APPLICATION_ID, PASSWORD)

	f, err := os.Open(FILEPATH)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	imageProcessingParams := &abbyysdk.ImageProcessingParams{
		ExportFormats: []abbyysdk.ExportFormat{abbyysdk.ExportFormatDocx, abbyysdk.ExportFormatTxt},
		Language:      "English,French",
	}

	taskInfo, err := ocrClient.ProcessImage(imageProcessingParams, f, FILEPATH)
	if err != nil {
		log.Fatal(err)
	}
	taskInfo, err = ocrClient.WaitForTask(taskInfo)
	if err != nil {
		log.Fatal(err)
	}

	for _, resultUrl := range taskInfo.ResultUrls {
		log.Println(resultUrl)
	}
}
