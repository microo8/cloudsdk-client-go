package abbyysdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/microo8/cloudsdk-client-go/models"
)

type OcrClient struct {
	client *HTTPClient
}

func NewOcrClient(host, applicationID, password string) *OcrClient {
	return &OcrClient{
		client: &HTTPClient{
			client:        &http.Client{},
			Host:          host,
			ApplicationID: applicationID,
			Password:      password,
		},
	}
}

func (c *OcrClient) startTask(
	httpMethod string,
	requestUrl string,
	params models.Params,
	fileStream io.Reader,
	fileName string,
	response interface{},
) error {
	resp, err := c.client.SendRequest(httpMethod, requestUrl, params, fileStream, fileName)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	d, _ := httputil.DumpResponse(resp, true)
	log.Println(string(d))
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return models.NewApiError(
			"Server responded with "+resp.Status+" status code.",
			resp.StatusCode,
			tryDeserializeError(responseData),
			resp.Header)
	}
	if err := json.Unmarshal(responseData, response); err != nil {
		return models.NewApiError(
			"Could not deserialize the response body: "+err.Error(),
			resp.StatusCode,
			models.NewErrorFromText(string(responseData)),
			resp.Header)
	}
	return nil
}

func (c *OcrClient) WaitForTask(taskInfo *models.TaskInfo) (*models.TaskInfo, error) {
	ti := taskInfo
	for {
		time.Sleep(time.Duration(ti.RequestStatusDelay) * time.Millisecond)
		ti, err := c.GetTaskStatus(taskInfo.TaskId)
		if err != nil {
			return nil, fmt.Errorf("getting task status: %w", err)
		}
		switch ti.Status {
		case models.TaskStatusInProgress, models.TaskStatusQueued, models.TaskStatusSubmitted:
		case models.TaskStatusDeleted:
			return nil, fmt.Errorf("task deleted")
		case models.TaskStatusProcessingFailed:
			return nil, fmt.Errorf("task processing failed: %w", ti.Error)
		case models.TaskStatusNotEnoughCredits:
			return nil, fmt.Errorf("not enough credits: %w", ti.Error)
		case models.TaskStatusCompleted:
			return ti, nil
		}
	}
}

func tryDeserializeError(responseData []byte) error {
	var e models.Error
	if err := json.Unmarshal(responseData, &e); err != nil {
		return models.NewErrorFromText(string(responseData))
	}
	return e
}

/**
 * The method loads the image, creates a processing task for the image with the specified parameters, and passes the task for processing.
 *
 * <b>Note:</b> This method allows you to specify up to three file formats for the result, in which case the server response
 * for the completed task will contain several result URLs. If there is not enough money on your account,
 * the task will be created, but will be suspended with {@link TaskStatus#NotEnoughCredits}
 * status. You can pass this task for processing using the {@link #processDocument(DocumentProcessingParams, bool)}
 * method after you have topped up your account. The task will not be created, if you exceed the limit of uploaded images.
 *
 * @param parameters Image processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>

func (c *OcrClient) ProcessImage(parameters *models.ImageProcessingParams, fileStream io.Reader, fileName string) (*models.TaskInfo, error) {
	resp := &models.TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessImageURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

/**
 * The method adds the image to the existing task or creates a new task for the image. This task is not passed for processing until
 * the {@link #processDocument(DocumentProcessingParams, bool)} or
 * {@link #processFields(FieldsProcessingParams, io.Reader, string, bool)} method is called for it.
 * Several images can be uploaded to one task
 * @param parameters Image submitting parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) SubmitImage(parameters *models.ImageSubmittingParams, fileStream io.Reader, fileName string) {
}

/**
 * The method starts the processing task with the specified parameters.
 *
 * <b>Note:</b> This method allows you to process several images using the same settings and obtain recognition
 * result as a multi-page document. You can upload several images to one task using {@link #submitImage(ImageSubmittingParams, io.Reader, string)} method.
 * It is also possible to specify up to three file formats for the result, in which case the server response for the completed
 * task will contain several result URLs. Only the task with {@link TaskStatus#Submitted},
 * {@link TaskStatus#Completed} or {@link TaskStatus#NotEnoughCredits}
 * status can be started using this method.
 *
 * @param parameters Document processing parameters
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessDocument(parameters models.DocumentProcessingParams) {
}

/**
 * The method allows you to recognize a business card on an image. The method loads the image, creates a processing task for the image
 * with the specified parameters, and passes the task for processing.
 * @param parameters Business card processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessBusinessCard(parameters models.BusinessCardProcessingParams, fileStream io.Reader, fileName string) {
}

/**
 * The method allows you to extract the value of a text field on an image. The method loads the image, creates a processing task for the image
 * with the specified parameters, and passes the task for processing. The result of recognition is returned in XML format.
 *
 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/quick-start-guide/text-fields-recognition/">How to Recognize Text Fields</a>
 *  to know how to tune the parameters.
 * <a href="https://www.ocrsdk.com/documentation/specifications/processing-profiles/"/>
 *
 * @param parameters Text field processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessTextField(parameters models.TextFieldProcessingParams, fileStream io.Reader, fileName string) {
}

/**
 * The method allows you to extract the value of a barcode on an image. The method loads the image, creates a processing task for
 * the image with the specified parameters, and passes the task for processing. The result of recognition is returned in XML format.
 * Binary data is returned in Base64 encoding.
 *
 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/quick-start-guide/barcode-ocr-sdk/">How to Recognize Barcodes</a>
 * to know another way of barcode recognition. {@link ProcessingProfile#BarcodeRecognition}
 * profile is used for processing. Information about processing profiles
 * <a href="https://www.ocrsdk.com/documentation/specifications/processing-profiles/"/>
 *
 * @param parameters Barcode field processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessBarcodeField(parameters models.BarcodeFieldProcessingParams, fileStream io.Reader, fileName string) {
}

/**
 * The method allows you to extract the value of a checkmark on an image. The method loads the image, creates a processing task for
 * the image with the specified parameters, and passes the task for processing. The result of recognition is returned in XML format.
 * The values of checkmarks are "checked", "unchecked", or "corrected".
 *
 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/processing-profiles/"/>
 *
 * @param parameters Checkmark field processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessCheckmarkField(parameters models.CheckmarkFieldProcessingParams, fileStream io.Reader, fileName string) {
}

/**
 * The method allows you to recognize several fields in a document. The method starts the processing task with the parameters of processing
 * specified in an XML file. Image files can be uploaded to the task using <see cref="SubmitImage"/> method. The result of recognition is
 * returned in XML format. Binary data is returned in Base64 encoding.
 *
 * <b>Note:</b> You can use the <a href="https://www.ocrsdk.com/schema/taskDescription-1.0.xsd">XSD schema</a>
 * of the XML file to create the file with necessary settings. See also the description of the tags and
 * several examples of XML files with settings in
 * <a href="https://www.ocrsdk.com/documentation/specifications/xml-scheme-field-settings/">XML Parameters of Field Recognition</a>.
 * Only the task with {@link TaskStatus#Submitted}, {@link TaskStatus#Completed}
 * or {@link TaskStatus#NotEnoughCredits} status can be started using this method.
 * Note that this method is most convenient when you process a large number of fields on one page: in this case the price of recognition
 * of all fields on one page does not exceed the price of recognition of a page of A4 size.
 *
 * @param parameters Fields processing parameters
 * @param fileStream XML File describing fields recognition settings
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessFields(parameters models.FieldsProcessingParams, fileStream io.Reader, fileName string) {
}

/**
 * This method finds a machine-readable zone on the image and extracts data from it. Machine-readable zone(MRZ) is typically found on
 * official travel or identity documents of many countries.It can have 2 lines or 3 lines of machine-readable data. This method allows to
 * process MRZ written in accordance with ICAO Document 9303 (endorsed by the International Organization for Standardization and the
 * International Electrotechnical Commission as ISO/IEC 7501-1)). The result of recognition is returned in XML format.
 *
 * <b>Note:</b> <a href="https://en.wikipedia.org/wiki/ICAO"/>
 * <a href="https://en.wikipedia.org/wiki/International_Organization_for_Standardization"/>
 * <a href="https://en.wikipedia.org/wiki/International_Electrotechnical_Commission"/>
 *
 * @param parameters Mrz processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessMrz(parameters models.MrzProcessingParams, fileStream io.Reader, fileName string) {
}

/**
 * Important: the technology fully supports the receipts issued in USA and France, other countries are currently supported in beta mode.
 * The method allows you to recognize the image of a receipt.The method loads the image, creates a processing task for the image with the
 * specified parameters, and passes the task for processing. The result is returned in XML format.
 *
 * <b>Note:</b> The elements and attributes of the resulting file are described in
 * <a href="https://en.wikipedia.org/wiki/ICAO"/>
 * For a step-by-step guide, see <a href="https://www.ocrsdk.com/documentation/quick-start-guide/receipt-recognition/">How to Recognize
 * Receipts.</a> The recommendations on preparing the input images can be found in
 * <a href="https://www.ocrsdk.com/documentation/hints-tips/photograph-scan-receipts/">Photographing and Scanning Receipts</a>.
 *
 * @param parameters Receipt processing parameters
 * @param fileStream Stream of the file with the image to recognize
 * @param fileName Name of the file with the image
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) ProcessReceipt(parameters models.ReceiptProccessingParams, fileStream io.Reader, fileName string) {
}

/**
 * The method returns the current status of the task and the URL of the result of processing for completed tasks.
 * Please note that the task status is not changed momentarily. Do not call this method more frequently than once in 2 or 3 seconds.
 * @param taskId Id of the task
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) GetTaskStatus(taskId string) (*models.TaskInfo, error) {
	resp := &models.TaskInfo{}
	if err := c.startTask(http.MethodGet, GetTaskStatusURL, &models.Task{TaskId: taskId}, nil, "", resp); err != nil {
		return nil, err
	}
	return resp, nil
}

/**
 * The method deletes the task and the images associated with this task from the ABBYY Cloud OCR SDK storage.
 * Only the tasks that have the status other than {@link TaskStatus#InProgress}
 * or {@link TaskStatus#Deleted} can be deleted.
 *
 * <b>Note:</b> If you try to delete the task that has already been deleted, the successful response is returned.
 * If you submit the same image to different tasks, to delete the image from the Cloud OCR SDK storage, you will need to call
 * this method for each task, which contains the image.
 *
 * @param taskId Id of the task
 * @return {@link TaskInfo}
 */
//TODO CompletableFuture<TaskInfo>
func (c *OcrClient) DeleteTask(taskId string) (*models.TaskInfo, error) {
	resp := &models.TaskInfo{}
	if err := c.startTask(http.MethodPost, DeleteTaskURL, &models.Task{TaskId: taskId}, nil, "", resp); err != nil {
		return nil, err
	}
	return resp, nil
}

/**
 * The method returns the list of tasks created by your application. By default, the {@link TaskStatus#Deleted}
 * tasks are included, but you can filter them out. This method may be useful if you need to compile an application usage
 * log for some period of time.
 *
 * <b>Note:</b> The tasks are ordered by the date of the last action with the task. This method can list up to 1000 tasks. If you process
 * a large number of tasks, call this method for shorter time periods.
 *
 * @param parameters Parameters for listing tasks
 * @return {@link TaskList}
 */
//TODO CompletableFuture<TaskList>
func (c *OcrClient) ListTasks(parameters models.TasksListingParams) {}

/**
 * The method returns the list of finished tasks. A task is finished if it has one of the following statuses:
 * {@link TaskStatus#Completed}, {@link TaskStatus#ProcessingFailed},
 * {@link TaskStatus#NotEnoughCredits}.
 *
 * <b>Note:</b> The tasks are ordered by the time of the end of processing. No more than 100 tasks can be returned at one method call. To
 * obtain more tasks, delete some finished tasks using the {@link #deleteTask(UUID)} method and then call this
 * method again.
 * The method may be useful if you work with a large number of tasks simultaneously. But there is no sense in calling this
 * method if your application does not have any incomplete tasks sent for the processing.
 * Please note that the task status is not changed momentarily. Do not call this method more frequently than once in 2 or
 *  seconds.
 * @return {@link TaskList}
 */
//TODO CompletableFuture<TaskList>
func (c *OcrClient) ListFinishedTask() {}

/**
 * This method allows you to receive information about the application type, its current balance and expiration date. You may
 * find it helpful for keeping the usage records.
 *
 * <b>Note:</b> The application is identified by its authentication information.
 * By default the call to this method is disabled for all applications. To enable getting the application information using
 * this method: 1) go to <a href="https://cloud.ocrsdk.com/"/> and log in 2) click Settings under your application's name
 * 3) on the next screen, click enable:
 * @return  {@link Application}
 */
//TODO CompletableFuture<Application>
func (c *OcrClient) GetApplicationInfo() {}
