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
	appID := os.Getenv("APPID")
	if appID == "" {
		log.Fatalf("must set APPID environment variable")
	}
	pass := os.Getenv("PASS")
	if pass == "" {
		log.Fatalf("must set PASS environment variable")
	}
}

func TestProcessImage(t *testing.T) {
	c := NewOcrClient(HOST, appID, pass)
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

	for _, resultUrl := range taskInfo.ResultUrls {
		log.Println(resultUrl)
	}

	/*
	   private static CompletableFuture<TaskInfo> processImageAsync() throws FileNotFoundException {
	       FileInputStream fileStream = new FileInputStream(FILEPATH);

	       ImageProcessingParams imageProcessingParams = new ImageProcessingParams();
	       imageProcessingParams.setExportFormats(new ExportFormat[]{ExportFormat.Docx, ExportFormat.Txt});
	       imageProcessingParams.setLanguage("English,French");

	       return ocrClient.processImageAsync(imageProcessingParams, fileStream, FILEPATH, true);
	*/
}
