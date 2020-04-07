package abbyysdk

import (
	"fmt"
	"net/http"
)

/*
Describes the API related error
*/
type ApiError struct {
	Message string
	/*
		HTTP Status code returned by a server
	*/
	StatusCode int

	/*
		Details about an error
	*/
	Err error

	/*
		Response headers
	*/
	Headers http.Header
}

/*
NewApiError instantiates the ApiException
message error message
statusCode HTTP Status code returned by a server
error Details about an error
headers Response headers
*/
func NewApiError(message string, statusCode int, error error, headers http.Header) ApiError {
	//return ApiError{message, statusCode, error, headers, nil}
	return ApiError{message, statusCode, error, headers}
}

func (err ApiError) Error() string {
	return err.Message + ": " + err.Err.Error()
}

//Error container
type Error struct {
	//Error object
	ErrorData ErrorData `json:"error"`
}

func NewErrorFromText(text string) Error {
	return Error{
		ErrorData: ErrorData{
			Message: text,
		},
	}
}

func (err Error) Error() string {
	return fmt.Sprintf("(%s) %s [%v]", err.ErrorData.Code, err.ErrorData.Message, err.ErrorData.Details)
}

//ErrorData describes the error details
type ErrorData struct {
	//The code of the error
	Code string

	//The message describing the error
	Message string

	//The description of the error occurence location
	Target string

	//Describes validation error
	Details []ErrorData
}
