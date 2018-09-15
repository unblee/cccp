package source

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type File struct {
	reader io.ReadCloser
	size   uint64
}

func NewFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read a file '%s'", path)
	}

	info, err := f.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get a file info '%s'", path)
	}

	return &File{
		reader: f,
		size:   uint64(info.Size()),
	}, nil
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *File) Close() error {
	return f.reader.Close()
}

func (f *File) Size() uint64 {
	return f.size
}
