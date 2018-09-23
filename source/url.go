package source

import (
	"net/http"

	"github.com/pkg/errors"
)

func NewURL(rawurl string) (*HTTPRequest, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new http request")
	}
	return NewHTTPRequest(req)
}
