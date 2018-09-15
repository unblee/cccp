package source

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type URL struct {
	reader io.ReadCloser
	size   uint64
}

func NewURL(rawurl string) (*URL, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to http get '%s'", rawurl)
	}

	size, err := strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "there is no content on '%s'", rawurl)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to get a file '%s': %d %s", rawurl, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return &URL{
		reader: resp.Body,
		size:   size,
	}, nil
}

func (u *URL) Read(p []byte) (n int, err error) {
	return u.reader.Read(p)
}

func (u *URL) Close() error {
	return u.reader.Close()
}

func (u *URL) Size() uint64 {
	return u.size
}
