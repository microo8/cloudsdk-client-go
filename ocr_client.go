//This library is intended to simplify the development of image recognition solutions with ABBYY Cloud OCR SDK API. This repo contains the client library enabling simple access to high-quality recognition technologies. For more information about the product visit https://www.ocrsdk.com/
package abbyysdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type httpClient struct {
	client        *http.Client
	Host          string
	ApplicationID string
	Password      string
}

func (client *httpClient) SendRequest(method, requestUrl string, params Params, fileStream io.Reader, fileName string) (*http.Response, error) {
	req, err := http.NewRequest(method, client.Host+"/"+requestUrl, fileStream)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	req.SetBasicAuth(client.ApplicationID, client.Password)
	query := req.URL.Query()
	for k, v := range params.Params() {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	return client.client.Do(req)
}

func (client *httpClient) DownloadResult(resultUrl string) (io.Reader, error) {
	req, err := http.NewRequest(http.MethodGet, resultUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type OcrClient struct {
	client *httpClient
}

func NewOcrClient(host, applicationID, password string) *OcrClient {
	return &OcrClient{
		client: &httpClient{
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
	params Params,
	fileStream io.Reader,
	fileName string,
	response interface{},
) error {
	resp, err := c.client.SendRequest(httpMethod, requestUrl, params, fileStream, fileName)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return NewApiError(
			"Server responded with "+resp.Status+" status code.",
			resp.StatusCode,
			tryDeserializeError(responseData),
			resp.Header)
	}
	if err := json.Unmarshal(responseData, response); err != nil {
		return NewApiError(
			"Could not deserialize the response body: "+err.Error(),
			resp.StatusCode,
			NewErrorFromText(string(responseData)),
			resp.Header)
	}
	return nil
}

func (c *OcrClient) WaitForTask(taskInfo *TaskInfo) (*TaskInfo, error) {
	ti := taskInfo
	for {
		time.Sleep(time.Duration(ti.RequestStatusDelay) * time.Millisecond)
		ti, err := c.GetTaskStatus(taskInfo.TaskId)
		if err != nil {
			return nil, fmt.Errorf("getting task status: %w", err)
		}
		switch ti.Status {
		case TaskStatusInProgress, TaskStatusQueued, TaskStatusSubmitted:
		case TaskStatusDeleted:
			return nil, fmt.Errorf("task deleted")
		case TaskStatusProcessingFailed:
			return nil, fmt.Errorf("task processing failed: %w", ti.Error)
		case TaskStatusNotEnoughCredits:
			return nil, fmt.Errorf("not enough credits: %w", ti.Error)
		case TaskStatusCompleted:
			return ti, nil
		}
	}
}

func (c *OcrClient) DownloadResults(taskInfo *TaskInfo) ([]io.Reader, error) {
	if taskInfo.Status != TaskStatusCompleted {
		return nil, fmt.Errorf("task is not completed")
	}
	if len(taskInfo.ResultUrls) == 0 {
		return nil, fmt.Errorf("task has no result urls")
	}
	var res []io.Reader
	for _, resultUrl := range taskInfo.ResultUrls {
		body, err := c.client.DownloadResult(resultUrl)
		if err != nil {
			return nil, err
		}
		res = append(res, body)
	}
	return res, nil
}

func tryDeserializeError(responseData []byte) error {
	var e Error
	if err := json.Unmarshal(responseData, &e); err != nil {
		return NewErrorFromText(string(responseData))
	}
	return e
}

//ProcessImage loads the image, creates a processing task for the image with the specified parameters, and passes the task for processing.
//Note: This method allows you to specify up to three file formats for the result, in which case the server response
//for the completed task will contain several result URLs. If there is not enough money on your account,
//the task will be created, but will be suspended with TaskStatusNotEnoughCredits
//status. You can pass this task for processing using the ProcessDocument(DocumentProcessingParams, bool)
//method after you have topped up your account. The task will not be created, if you exceed the limit of uploaded images.
//parameters Image processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessImage(parameters *ImageProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessImageURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//SubmitImage adds the image to the existing task or creates a new task for the image. This task is not passed for processing until
//the ProcessDocument(DocumentProcessingParams, bool) or
//ProcessFields(FieldsProcessingParams, io.Reader, string, bool) method is called for it.
//Several images can be uploaded to one task
//parameters Image submitting parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) SubmitImage(parameters *ImageSubmittingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, SubmitImageURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessDocument starts the processing task with the specified parameters.
//Note: This method allows you to process several images using the same settings and obtain recognition
//result as a multi-page document. You can upload several images to one task using SubmitImage(ImageSubmittingParams, io.Reader, string) method.
//It is also possible to specify up to three file formats for the result, in which case the server response for the completed
//task will contain several result URLs. Only the task with TaskStatusSubmitted,
//TaskStatusCompleted or TaskStatusNotEnoughCredits
//status can be started using this method.
//parameters Document processing parameters
func (c *OcrClient) ProcessDocument(parameters *DocumentProcessingParams) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessDocumentURL, parameters, nil, "", resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessBusinessCard allows you to recognize a business card on an image. The method loads the image, creates a processing task for the image
//with the specified parameters, and passes the task for processing.
//parameters Business card processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessBusinessCard(parameters *BusinessCardProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessBusinessCardURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessTextField allows you to extract the value of a text field on an image. The method loads the image, creates a processing task for the image
//with the specified parameters, and passes the task for processing. The result of recognition is returned in XML format.
//Note: See https://www.ocrsdk.com/documentation/quick-start-guide/text-fields-recognition How to Recognize Text Fields
//to know how to tune the parameters.
//https://www.ocrsdk.com/documentation/specifications/processing-profiles
//parameters Text field processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessTextField(parameters *TextFieldProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessTextFieldURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessBarcodeField allows you to extract the value of a barcode on an image. The method loads the image, creates a processing task for
//the image with the specified parameters, and passes the task for processing. The result of recognition is returned in XML format.
//Binary data is returned in Base64 encoding.
//Note: See https://www.ocrsdk.com/documentation/quick-start-guide/barcode-ocr-sdk
//How to Recognize Barcodes to know another way of barcode recognition. ProcessingProfileBarcodeRecognition
//profile is used for processing. Information about processing profiles
//https://www.ocrsdk.com/documentation/specifications/processing-profiles
//parameters Barcode field processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessBarcodeField(parameters *BarcodeFieldProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessBarcodeFieldURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessCheckmarkField allows you to extract the value of a checkmark on an image. The method loads the image, creates a processing task for
//the image with the specified parameters, and passes the task for processing. The result of recognition is returned in XML format.
//The values of checkmarks are "checked", "unchecked", or "corrected".
//Note: See https://www.ocrsdk.com/documentation/specifications/processing-profiles
//parameters Checkmark field processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessCheckmarkField(parameters *CheckmarkFieldProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessCheckmarkFieldURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessFields allows you to recognize several fields in a document. The method starts the processing task with the parameters of processing
//specified in an XML file. Image files can be uploaded to the task using SubmitImage method. The result of recognition is
//returned in XML format. Binary data is returned in Base64 encoding.
//Note: You can use the https://www.ocrsdk.com/schema/taskDescription-1.0.xsd
//XSD schema of the XML file to create the file with necessary settings. See also the description of the tags and
//several examples of XML files with settings in
//https://www.ocrsdk.com/documentation/specifications/xml-scheme-field-settings
//XML Parameters of Field Recognition.
//Only the task with TaskStatusSubmitted, TaskStatusCompleted
//or TaskStatusNotEnoughCredits status can be started using this method.
//Note that this method is most convenient when you process a large number of fields on one page: in this case the price of recognition
//of all fields on one page does not exceed the price of recognition of a page of A4 size.
//parameters Fields processing parameters
//fileStream XML File describing fields recognition settings
//fileName Name of the file with the image

func (c *OcrClient) ProcessFields(parameters *FieldsProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessFieldsURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ProcessMrz finds a machine-readable zone on the image and extracts data from it. Machine-readable zone(MRZ) is typically found on
//official travel or identity documents of many countries.It can have 2 lines or 3 lines of machine-readable data. This method allows to
//process MRZ written in accordance with ICAO Document 9303 (endorsed by the International Organization for Standardization and the
//International Electrotechnical Commission as ISO/IEC 7501-1)). The result of recognition is returned in XML format.
//
//Note: https://en.wikipedia.org/wiki/ICAO
//
//https://en.wikipedia.org/wiki/International_Organization_for_Standardization
//
//https://en.wikipedia.org/wiki/International_Electrotechnical_Commission
//
//parameters Mrz processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessMrz(parameters *MrzProcessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessMrzURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//Important: the technology fully supports the receipts issued in USA and France, other countries are currently supported in beta mode.
//ProcessReceipt allows you to recognize the image of a receipt.The method loads the image, creates a processing task for the image with the
//specified parameters, and passes the task for processing. The result is returned in XML format.
//Note: The elements and attributes of the resulting file are described in
//https://en.wikipedia.org/wiki/ICAO
//For a step-by-step guide, see https://www.ocrsdk.com/documentation/quick-start-guide/receipt-recognition
//How to Recognize Receipts. The recommendations on preparing the input images can be found in
//https://www.ocrsdk.com/documentation/hints-tips/photograph-scan-receipts
//Photographing and Scanning Receipts
//parameters Receipt processing parameters
//fileStream Stream of the file with the image to recognize
//fileName Name of the file with the image
func (c *OcrClient) ProcessReceipt(parameters *ReceiptProccessingParams, fileStream io.Reader, fileName string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, ProcessReceiptURL, parameters, fileStream, fileName, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//GetTaskStatus returns the current status of the task and the URL of the result of processing for completed tasks.
//Please note that the task status is not changed momentarily. Do not call this method more frequently than once in 2 or 3 seconds.
//taskId Id of the task
func (c *OcrClient) GetTaskStatus(taskId string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodGet, GetTaskStatusURL, &Task{TaskId: taskId}, nil, "", resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//DeleteTask deletes the task and the images associated with this task from the ABBYY Cloud OCR SDK storage.
//Only the tasks that have the status other than TaskStatusInProgress
//or TaskStatusDeleted can be deleted.
//Note: If you try to delete the task that has already been deleted, the successful response is returned.
//If you submit the same image to different tasks, to delete the image from the Cloud OCR SDK storage, you will need to call
//this method for each task, which contains the image.
//taskId Id of the task
func (c *OcrClient) DeleteTask(taskId string) (*TaskInfo, error) {
	resp := &TaskInfo{}
	if err := c.startTask(http.MethodPost, DeleteTaskURL, &Task{TaskId: taskId}, nil, "", resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//ListTasks returns the list of tasks created by your application. By default, the TaskStatusDeleted
//tasks are included, but you can filter them out. This method may be useful if you need to compile an application usage
//log for some period of time.
//Note: The tasks are ordered by the date of the last action with the task. This method can list up to 1000 tasks. If you process
//a large number of tasks, call this method for shorter time periods.
//parameters Parameters for listing tasks
func (c *OcrClient) ListTasks(parameters *TasksListingParams) {}

//ListFinishedTask returns the list of finished tasks. A task is finished if it has one of the following statuses:
//TaskStatusCompleted, TaskStatusProcessingFailed},
//TaskStatusNotEnoughCredits}.
//Note: The tasks are ordered by the time of the end of processing. No more than 100 tasks can be returned at one method call. To
//obtain more tasks, delete some finished tasks using the deleteTask(UUID) method and then call this
//method again.
//The method may be useful if you work with a large number of tasks simultaneously. But there is no sense in calling this
//method if your application does not have any incomplete tasks sent for the processing.
//Please note that the task status is not changed momentarily. Do not call this method more frequently than once in 2 or
//seconds.
func (c *OcrClient) ListFinishedTask() {}

//GetApplicationInfo allows you to receive information about the application type, its current balance and expiration date. You may
//find it helpful for keeping the usage records.
//Note: The application is identified by its authentication information.
//By default the call to this method is disabled for all applications. To enable getting the application information using
//this method:
//1) go to https://cloud.ocrsdk.com and log in
//2) click Settings under your application's name
//3) on the next screen, click enable:
func (c *OcrClient) GetApplicationInfo() {}
