package abbyysdk

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/microo8/cloudsdk-client-go/models"
)

type HTTPClient struct {
	client        *http.Client
	Host          string
	ApplicationID string
	Password      string
}

func (client *HTTPClient) SendRequest(method, requestUrl string, params models.Params, fileStream io.Reader, fileName string) (*http.Response, error) {
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
	d, _ := httputil.DumpRequest(req, false)
	log.Println(string(d))
	return client.client.Do(req)
}
