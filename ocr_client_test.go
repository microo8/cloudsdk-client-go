package abbyysdk

import (
	"log"
	"os"
	"testing"

	"github.com/microo8/cloudsdk-client-go/models"
)

const (
	HOST = "https://cloud-eu.ocrsdk.com"
)

var (
	appID string
	pass  string
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	appID = os.Getenv("APPID")
	if appID == "" {
		log.Fatalf("must set APPID environment variable")
	}
	pass = os.Getenv("PASS")
	if pass == "" {
		log.Fatalf("must set PASS environment variable")
	}
}

func TestProcessImage(t *testing.T) {
	c := NewOcrClient(HOST, appID, pass)
	f, err := os.Open("resources/processImage.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	taskInfo, err := c.ProcessImage(&models.ImageProcessingParams{
		ExportFormats: []models.ExportFormat{models.ExportFormatDocx, models.ExportFormatTxt},
		Language:      "English,French",
	},
		f,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	taskInfo, err = c.WaitForTask(taskInfo)
	if err != nil {
		t.Fatal(err)
	}
	for _, resultUrl := range taskInfo.ResultUrls {
		log.Println(resultUrl)
	}
}
