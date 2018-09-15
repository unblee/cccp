package source

import (
	"io"
)

type Source interface {
	io.ReadCloser
	Size() uint64
}
