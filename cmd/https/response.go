package https

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type HTTPResponse struct {
	Body       string
	Bytes      []byte
	Status     string
	StatusCode int
}

func (res *HTTPResponse) extractBody() {
	if len(res.Bytes) > 0 {
		res.Body = string(res.Bytes)
	}
}

func (res *HTTPResponse) extractResponseDetails(resp *http.Response) {
	if resp != nil {
		res.extractBody()
		res.Status = resp.Status
		res.StatusCode = resp.StatusCode
	}
}

func (res *HTTPResponse) extractRawBytes(resp *http.Response) {
	var err error
	if resp.Body != nil {
		defer func() {
			if err = resp.Body.Close(); err != nil {
				log.Println(fmt.Sprintf("error encountered extracting bytes from response: %s", err))
			}
		}()
		if res.Bytes, err = io.ReadAll(resp.Body); err != nil {
			log.Println(fmt.Sprintf("unable to retrieve bytes from response: %s", err))
		}
	}

}

func newResponse() HTTPResponse {
	return HTTPResponse{
		Status:     "Internal Server Error",
		StatusCode: 501,
		Body:       "",
	}
}
