package abbyysdk

import (
	"fmt"
	"io"
	"net/http"

	"github.com/microo8/cloudsdk-client-go/models"
)

type HTTPClient struct {
	client        *http.Client
	Host          string
	ApplicationID string
	Password      string
}

func (client *HTTPClient) SendRequest(requestUrl string, params models.Params, fileStream io.Reader, fileName string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, requestUrl, fileStream)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	query := req.URL.Query()
	for k, v := range params.Params() {
		query.Add(k, v)
	}
	return client.client.Do(req)
}
