package source

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type HTTPRequest struct {
	reader io.ReadCloser
	size   uint64
}

func NewHTTPRequest(req *http.Request) (*HTTPRequest, error) {
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to http %s '%s'", req.Method, req.RequestURI)
	}

	size, err := strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "there is no content on '%s'", req.RequestURI)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to get response body '%s': %d %s", req.RequestURI, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return &HTTPRequest{
		reader: resp.Body,
		size:   size,
	}, nil
}

func (req *HTTPRequest) Read(p []byte) (n int, err error) {
	return req.reader.Read(p)
}

func (req *HTTPRequest) Close() error {
	return req.reader.Close()
}

func (req *HTTPRequest) Size() uint64 {
	return req.size
}
