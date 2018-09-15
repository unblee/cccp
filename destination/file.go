package destination

import (
	"io"
	"os"
	"sync/atomic"

	"github.com/pkg/errors"
)

type File struct {
	writer  io.WriteCloser
	written uint64
}

func NewFile(path string) (*File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a file '%s'", path)
	}

	return &File{
		writer: f,
	}, nil
}

func (f *File) Write(p []byte) (n int, err error) {
	atomic.AddUint64(&f.written, uint64(len(p)))
	return f.writer.Write(p)
}

func (f *File) Close() error {
	return f.writer.Close()
}

func (f *File) Written() uint64 {
	return atomic.LoadUint64(&f.written)
}
