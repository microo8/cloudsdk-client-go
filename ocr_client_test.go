package abbyysdk

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/microo8/cloudsdk-client-go/abbyyxml"
)

const (
	HOST = "https://cloud-eu.ocrsdk.com"

	prefix             = "resources/"
	IMAGE              = prefix + "processImage.jpg"
	BUSINESS_CARD      = prefix + "processBusinessCard.jpg"
	FIELDS_XML_CONFIG  = prefix + "ProcessFieldsXmlConfig.xml"
	MRZ                = prefix + "processMrz.jpg"
	FIELDS             = prefix + "processFields.tif"
	UNSUPPORTED_FORMAT = prefix + "unsupported.properties"
)

var (
	appID string
	pass  string
	c     *OcrClient
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
	c = NewOcrClient(HOST, appID, pass)
}

func submitImage(filename, taskId string) (*TaskInfo, error) {
	fileStream, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileStream.Close()

	imageSubmittingParams := &ImageSubmittingParams{}
	imageSubmittingParams.TaskId = taskId

	submitImageTask, err := c.SubmitImage(imageSubmittingParams, fileStream, filename)
	if err != nil {
		return nil, fmt.Errorf("submitting image: %w", err)
	}
	if submitImageTask == nil {
		return nil, fmt.Errorf("submitImageTask is nil")
	}
	if submitImageTask.TaskId == "" {
		return nil, fmt.Errorf("submitImageTask.TaskId is empty")
	}
	if submitImageTask.Status != TaskStatusSubmitted {
		return nil, fmt.Errorf("submitImageTask.Status is not submitted: %s", submitImageTask.Status)
	}

	return submitImageTask, nil
}

func checkResultTask(taskInfo *TaskInfo, taskId string, resultUrls int, taskStatus TaskStatus) error {
	if taskInfo == nil {
		return fmt.Errorf("taskInfo is nil")
	}
	if taskId != "" {
		if taskId != taskInfo.TaskId {
			return fmt.Errorf("taskInfo.TaskId is not equal")
		}
	}
	if taskStatus != taskInfo.Status {
		return fmt.Errorf("taskInfo.Status is not equal")
	}
	if resultUrls != len(taskInfo.ResultUrls) {
		return fmt.Errorf("taskInfo.ResultUrls len is not equal")
	}
	return nil
}

func TestProcessImage(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(IMAGE)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()
	imageProcessingParams := &ImageProcessingParams{
		Language:      "English",
		ExportFormats: []ExportFormat{ExportFormatDocx},
	}
	processImageTask, err := c.ProcessImage(
		imageProcessingParams,
		fileStream,
		IMAGE,
	)
	if err != nil {
		t.Fatal(err)
	}
	processImageTask, err = c.WaitForTask(processImageTask)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkResultTask(processImageTask, "", 1, TaskStatusCompleted); err != nil {
		t.Fatal(err)
	}
}

func TestSubmitImage(t *testing.T) {
	t.Parallel()
	first, err := submitImage(IMAGE, "")
	if err != nil {
		t.Fatal(err)
	}
	second, err := submitImage(FIELDS, first.TaskId)
	if err != nil {
		t.Fatal(err)
	}
	if 1 != first.FilesCount {
		t.Errorf("submitImage wrong files count: %d", first.FilesCount)
	}
	if 2 != second.FilesCount {
		t.Errorf("submitImage wrong files count: %d", second.FilesCount)
	}
}

func TestProcessDocument(t *testing.T) {
	t.Parallel()
	submitImageTask, err := submitImage(IMAGE, "")
	if err != nil {
		t.Fatal(err)
	}

	documentProcessingParams := &DocumentProcessingParams{
		Language:      "English",
		Profile:       ProcessingProfileDocumentConversion,
		ExportFormats: []ExportFormat{ExportFormatPdfSearchable, ExportFormatRtf},
	}
	documentProcessingParams.TaskId = submitImageTask.TaskId

	processDocumentTask, err := c.ProcessDocument(documentProcessingParams)
	if err != nil {
		t.Fatal(err)
	}
	processDocumentTask, err = c.WaitForTask(processDocumentTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processDocumentTask, submitImageTask.TaskId, 2, TaskStatusCompleted)
}

func TestProcessBusinessCard(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(BUSINESS_CARD)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	businessCardProcessingParams := &BusinessCardProcessingParams{
		Language:     "English",
		ExportFormat: BusinessCardExportFormatXml,
	}

	processBusinessCardTask, err := c.ProcessBusinessCard(
		businessCardProcessingParams,
		fileStream,
		BUSINESS_CARD,
	)
	if err != nil {
		t.Fatal(err)
	}
	processBusinessCardTask, err = c.WaitForTask(processBusinessCardTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processBusinessCardTask, "", 1, TaskStatusCompleted)
}

func TestProcessTextField(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(BUSINESS_CARD)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	textFieldProcessingParams := &TextFieldProcessingParams{
		Language: "English",
		Region:   "140,550,1130,700",
	}

	processTextFieldTask, err := c.ProcessTextField(
		textFieldProcessingParams,
		fileStream,
		BUSINESS_CARD,
	)
	if err != nil {
		t.Fatal(err)
	}
	processTextFieldTask, err = c.WaitForTask(processTextFieldTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processTextFieldTask, "", 1, TaskStatusCompleted)
}

func TestProcessBarcodeField(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(FIELDS)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	barcodeFieldProcessingParams := &BarcodeFieldProcessingParams{
		BarcodeTypes: []BarcodeType{BarcodeTypeEan8},
		Region:       "1800,3200,2250,3400",
	}

	processBarcodeFieldTask, err := c.ProcessBarcodeField(
		barcodeFieldProcessingParams,
		fileStream,
		FIELDS,
	)
	if err != nil {
		t.Fatal(err)
	}
	processBarcodeFieldTask, err = c.WaitForTask(processBarcodeFieldTask)
	if err != nil {
		t.Fatal(err)
	}
	results, err := c.DownloadResults(processBarcodeFieldTask)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		var doc abbyyxml.Document
		if err := xml.NewDecoder(r).Decode(&doc); err != nil {
			t.Fatal(err)
		}
		log.Println(doc)
	}
	checkResultTask(processBarcodeFieldTask, "", 1, TaskStatusCompleted)
}

func TestProcessCheckmarkField(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(FIELDS)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	checkmarkFieldProcessingParams := &CheckmarkFieldProcessingParams{
		CheckmarkType: CheckmarkTypeSquare,
		Region:        "700,930,800,1030",
	}

	processCheckmarkFieldTask, err := c.ProcessCheckmarkField(
		checkmarkFieldProcessingParams,
		fileStream,
		FIELDS,
	)
	if err != nil {
		t.Fatal(err)
	}
	processCheckmarkFieldTask, err = c.WaitForTask(processCheckmarkFieldTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processCheckmarkFieldTask, "", 1, TaskStatusCompleted)
}

func TestProcessFields(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(FIELDS_XML_CONFIG)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	submitImageTask, err := submitImage(FIELDS, "")
	if err != nil {
		t.Fatal(err)
	}

	b := true
	fieldProcessingParams := &FieldsProcessingParams{
		WriteRecognitionVariants: &b,
	}
	fieldProcessingParams.TaskId = submitImageTask.TaskId

	processFieldsTask, err := c.ProcessFields(fieldProcessingParams, fileStream, FIELDS_XML_CONFIG)
	if err != nil {
		t.Fatal(err)
	}
	processFieldsTask, err = c.WaitForTask(processFieldsTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processFieldsTask, submitImageTask.TaskId, 1, TaskStatusCompleted)
}

func TestProcessMrz(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(MRZ)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	mrzProcessingParams := &MrzProcessingParams{
		Description: "Task description: MRZ processing.",
		PdfPassword: "test",
	}

	processMrzTask, err := c.ProcessMrz(
		mrzProcessingParams,
		fileStream,
		FIELDS,
	)
	if err != nil {
		t.Fatal(err)
	}
	processMrzTask, err = c.WaitForTask(processMrzTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processMrzTask, "", 1, TaskStatusCompleted)
}

func TestProcessReceipt(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(IMAGE)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	receiptProccessingParams := &ReceiptProccessingParams{
		Countries: []ReceiptRecognizingCountry{ReceiptRecognizingCountryRussia},
	}

	processReceiptTask, err := c.ProcessReceipt(
		receiptProccessingParams,
		fileStream,
		FIELDS,
	)

	if err != nil {
		t.Fatal(err)
	}
	processReceiptTask, err = c.WaitForTask(processReceiptTask)
	if err != nil {
		t.Fatal(err)
	}

	checkResultTask(processReceiptTask, "", 1, TaskStatusCompleted)
}

func TestGetTaskStatus(t *testing.T) {
	t.Parallel()
	submitImageTask, err := submitImage(IMAGE, "")
	if err != nil {
		t.Fatal(err)
	}
	resultTask, err := c.GetTaskStatus(submitImageTask.TaskId)
	if err != nil {
		t.Fatal(err)
	}

	if resultTask == nil {
		t.Fatal("resultTask is nil")
	}
	if submitImageTask.TaskId != resultTask.TaskId {
		t.Fatal("resultTask.TaskId not equal")
	}
	if TaskStatusSubmitted != resultTask.Status {
		t.Fatal("wrong status")
	}
}

func TestDeleteTask(t *testing.T) {
	t.Parallel()
	submitImageTask, err := submitImage(IMAGE, "")
	if err != nil {
		t.Fatal(err)
	}
	deletedTask, err := c.DeleteTask(submitImageTask.TaskId)
	if err != nil {
		t.Fatal(err)
	}
	resultTask, err := c.GetTaskStatus(deletedTask.TaskId)
	if err != nil {
		t.Fatal(err)
	}
	if deletedTask == nil {
		t.Fatal("deletedTask is nil")
	}
	if resultTask == nil {
		t.Fatal("resultTask is nil")
	}

	if submitImageTask.TaskId != deletedTask.TaskId {
		t.Fatal("ids not equal")
	}
	if deletedTask.TaskId != resultTask.TaskId {
		t.Fatal("ids not equal")
	}

	if TaskStatusSubmitted != submitImageTask.Status {
		t.Fatal("status not submitted")
	}
	if TaskStatusDeleted != deletedTask.Status {
		t.Fatal("status not deleted")
	}
	if TaskStatusDeleted != resultTask.Status {
		t.Fatal("status not deleted")
	}
}

func TestListTasks(t *testing.T) {
	t.Parallel()
	submitImageTask, err := submitImage(IMAGE, "")
	if err != nil {
		t.Fatal(err)
	}

	b := true
	tasksListingParams := &TasksListingParams{
		ExcludeDeleted: &b,
	}

	taskList, err := c.ListTasks(tasksListingParams)
	if err != nil {
		t.Fatal(err)
	}
	if taskList == nil {
		t.Fatal("taskList is nil")
	}
	if len(taskList) == 0 {
		t.Fatal("taskList is empty")
	}
	if !idIn(submitImageTask.TaskId, taskList) {
		t.Fatal("submitImageTask.TaskId not in taskList")
	}
}

func idIn(id string, taskList []*TaskInfo) bool {
	for _, task := range taskList {
		if task.TaskId == id {
			return true
		}
	}
	return false
}

func TestListFinishedTasks(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(IMAGE)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	submitImageTask, err := submitImage(IMAGE, "")
	if err != nil {
		t.Fatal(err)
	}

	processImageTask, err := c.ProcessImage(
		&ImageProcessingParams{},
		fileStream,
		IMAGE,
	)
	if err != nil {
		t.Fatal(err)
	}
	processImageTask, err = c.WaitForTask(processImageTask)
	if err != nil {
		t.Fatal(err)
	}

	listFinishedTasks, err := c.ListFinishedTask()
	if err != nil {
		t.Fatal(err)
	}
	if listFinishedTasks == nil {
		t.Fatal("listFinishedTasks is nil")
	}
	if len(listFinishedTasks) == 0 {
		t.Fatal("listFinishedTasks len is 0")
	}
	if len(listFinishedTasks) > 100 {
		t.Fatal("listFinishedTasks len more than 100")
	}

	if idIn(submitImageTask.TaskId, listFinishedTasks) {
		t.Fatal("submitImageTask.TaskId in listFinishedTasks")
	}
	if !idIn(processImageTask.TaskId, listFinishedTasks) {
		t.Fatal("processImageTask.TaskId not in listFinishedTasks")
	}

	b := true
	tasksListingParams := &TasksListingParams{
		ExcludeDeleted: &b,
	}

	taskList, err := c.ListTasks(tasksListingParams)
	if err != nil {
		t.Fatal(err)
	}
	if taskList == nil {
		t.Fatal("taskList is nil")
	}
	if len(taskList) == 0 {
		t.Fatal("taskList is empty")
	}
	if !idIn(submitImageTask.TaskId, taskList) {
		t.Fatal("submitImageTask.TaskId not in taskList")
	}

}

func TestGetApplicationInfo(t *testing.T) {
	t.Parallel()
	application, err := c.GetApplicationInfo()
	if err != nil {
		t.Fatal(err)
	}
	if application == nil {
		t.Fatal("application is nil")
	}
	if application.ID != appID {
		t.Fatalf("application ID not equal: %s - %s", application.ID, appID)
	}
}

func TestInvalidTaskId(t *testing.T) {
	t.Parallel()
	invalidTaskId := "163024a1-eee6-42ea-a577-a4b936dcb250"
	if _, err := c.GetTaskStatus(invalidTaskId); err == nil {
		t.Fatal("ApiError expected")
	}
}

func TestNullTaskId(t *testing.T) {
	t.Parallel()
	if _, err := c.GetTaskStatus(""); err == nil {
		t.Fatal("ApiError expected")
	}
}

func TestInvalidRegion(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(FIELDS)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	barcodeFieldProcessingParams := &BarcodeFieldProcessingParams{
		Region: "some_invalid_region",
	}

	if _, err := c.ProcessBarcodeField(
		barcodeFieldProcessingParams,
		fileStream,
		FIELDS,
	); err == nil {
		t.Fatal("ApiError expected")
	}
}

func TestUnsupportedFormat(t *testing.T) {
	t.Parallel()
	fileStream, err := os.Open(UNSUPPORTED_FORMAT)
	if err != nil {
		t.Fatal(err)
	}
	defer fileStream.Close()

	barcodeFieldProcessingParams := &BarcodeFieldProcessingParams{}

	if _, err := c.ProcessBarcodeField(
		barcodeFieldProcessingParams,
		fileStream,
		UNSUPPORTED_FORMAT,
	); err == nil {
		t.Fatal("ApiError expected")
	}
}

func TestNoImage(t *testing.T) {
	t.Parallel()
	barcodeFieldProcessingParams := &BarcodeFieldProcessingParams{}

	if _, err := c.ProcessBarcodeField(
		barcodeFieldProcessingParams,
		nil,
		"",
	); err == nil {
		t.Fatal("ApiError expected")
	}
}
