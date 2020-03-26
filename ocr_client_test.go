package abbyysdk

import (
	"fmt"
	"log"
	"os"
	"testing"
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

/*


   func TestGetTaskStatus(t *testing.T) {
	t.Parallel()
       TaskInfo submitImageTask = submitImage(.IMAGE);
       TaskInfo resultTask = c.getTaskStatus(submitImageTask.getTaskId()).get();

       Assert.assertNotNull(resultTask);
       Assert.assertEquals(submitImageTask.getTaskId(), resultTask.getTaskId());
       Assert.assertEquals(TaskStatus.Submitted, resultTask.getStatus());
   }


   func TestDeleteTask(t *testing.T) {
	t.Parallel()
       TaskInfo submitImageTask = submitImage(.IMAGE);
       TaskInfo deletedTask = c.deleteTask(submitImageTask.getTaskId()).get();
       TaskInfo resultTask = c.getTaskStatus(deletedTask.getTaskId()).get();

       Assert.assertNotNull(deletedTask);
       Assert.assertNotNull(resultTask);

       Assert.assertEquals(submitImageTask.getTaskId(), deletedTask.getTaskId());
       Assert.assertEquals(deletedTask.getTaskId(), resultTask.getTaskId());

       Assert.assertEquals(TaskStatus.Submitted, submitImageTask.getStatus());
       Assert.assertEquals(TaskStatus.Deleted, deletedTask.getStatus());
       Assert.assertEquals(TaskStatus.Deleted, resultTask.getStatus());
   }


   func TestListTasks(t *testing.T) {
	t.Parallel()
       TaskInfo submitImageTask = submitImage(.IMAGE);

       TasksListingParams tasksListingParams = new TasksListingParams();
       tasksListingParams.setExcludeDeleted(true);

       TaskList taskList = c.listTasks(tasksListingParams).get();

       Assert.assertNotNull(taskList);
       Assert.assertTrue(taskList.getTasks().size() > 0);
       Assert.assertTrue(taskList.getTasks().stream()
               .anyMatch(task -> task.getTaskId().equals(submitImageTask.getTaskId())));


   }


   func TestListFinishedTasks(t *testing.T) {
	t.Parallel()
       FileInputStream fileStream = new FileInputStream(.IMAGE);

       TaskInfo submitImageTask = submitImage(.IMAGE);

       TaskInfo processImageTask = c.processImage(
               new ImageProcessingParams(),
               fileStream,
               .IMAGE,
               true
       ).get();

       TaskList listFinishedTasks = c.listFinishedTask().get();

       Assert.assertNotNull(listFinishedTasks);

       Assert.assertTrue(listFinishedTasks.getTasks().size() > 0);
       Assert.assertTrue(listFinishedTasks.getTasks().size() <= 100);

       Assert.assertFalse(listFinishedTasks.getTasks().stream().anyMatch(task -> task.getTaskId().equals(submitImageTask.getTaskId())));
       Assert.assertTrue(listFinishedTasks.getTasks().stream().anyMatch(task -> task.getTaskId().equals(processImageTask.getTaskId())));

       TasksListingParams tasksListingParams = new TasksListingParams();
       tasksListingParams.setExcludeDeleted(true);

       TaskList taskList = c.listTasks(tasksListingParams).get();

       Assert.assertNotNull(taskList);
       Assert.assertTrue(taskList.getTasks().size() > 0);
       Assert.assertTrue(taskList.getTasks().stream()
               .anyMatch(task -> task.getTaskId().equals(submitImageTask.getTaskId())));

   }


   func TestGetApplicationInfo(t *testing.T) {
	t.Parallel()
       Application application = c.getApplicationInfo().get();

       Assert.assertNotNull(application);
       Assert.assertEquals(testConfig.getApplicationId(), application.getId());
   }


   func TestInvalidTaskId(t *testing.T) {
	t.Parallel()
       // Due to Java async-call mechanism, any exception, thrown during asynchronous method is wrapped with ExecutionException
       // So, to obtain original exception, one has to get cause of ExecutionException instance
       try {
           try {
               UUID invalidTaskId = UUID.randomUUID();
               c.getTaskStatus(invalidTaskId).get();

               Assert.fail("ApiException expected");
           } catch (ExecutionException e) {
               throw e.getCause();
           }
       } catch (ApiException e) {
           ShouldHelper.checkException(e, ErrorCode.InvalidArgument, ValidationErrorCode.InvalidParameterValue, ErrorTarget.TaskId);
       }
   }


   func TestNullTaskId(t *testing.T) {
	t.Parallel()
       try {
           try {
               c.getTaskStatus(null).get();

               Assert.fail("ApiException expected");
           } catch (ExecutionException e) {
               throw e.getCause();
           }
       } catch (ApiException e) {
           ShouldHelper.checkException(e, ErrorCode.InvalidArgument, ValidationErrorCode.MissingArgument, ErrorTarget.TaskId);
       }
   }


   func TestInvalidRegion(t *testing.T) {
	t.Parallel()
       try {
           try {
               FileInputStream fileStream = new FileInputStream(.FIELDS);

               BarcodeFieldProcessingParams barcodeFieldProcessingParams = new BarcodeFieldProcessingParams();
               barcodeFieldProcessingParams.setRegion("some_invalid_region");

               c.processBarcodeField(
                       barcodeFieldProcessingParams,
                       fileStream,
                       .FIELDS,
                       true
               ).get();

               Assert.fail("ApiException expected");
           } catch (ExecutionException e) {
               throw e.getCause();
           }
       } catch (ApiException e) {
           ShouldHelper.checkException(e, ErrorCode.InvalidArgument, ValidationErrorCode.InvalidParameterValue, ErrorTarget.Region);
       }
   }


   func TestUnsupportedFormat(t *testing.T) {
	t.Parallel()
       try {
           try {
               FileInputStream fileStream = new FileInputStream(.UNSUPPORTED_FORMAT);

               BarcodeFieldProcessingParams barcodeFieldProcessingParams = new BarcodeFieldProcessingParams();

               c.processBarcodeField(
                       barcodeFieldProcessingParams,
                       fileStream,
                       .UNSUPPORTED_FORMAT,
                       true
               ).get();

               Assert.fail("ApiException expected");
           } catch (ExecutionException e) {
               throw e.getCause();
           }
       } catch (ApiException e) {
           ShouldHelper.checkException(e, ErrorCode.FileFormatUnsupported);
       }
   }


   func TestNoImage(t *testing.T) {
	t.Parallel()
       try {
           try {
               BarcodeFieldProcessingParams barcodeFieldProcessingParams = new BarcodeFieldProcessingParams();

               c.processBarcodeField(
                       barcodeFieldProcessingParams,
                       null,
                       null,
                       true
               ).get();

               Assert.fail("ApiException expected");
           } catch (ExecutionException e) {
               throw e.getCause();
           }
       } catch (ApiException e) {
           ShouldHelper.checkException(e, ErrorCode.InvalidArgument, ValidationErrorCode.MissingArgument, ErrorTarget.File);
       }
   }

*/
