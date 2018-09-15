package cccp

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/unblee/cccp/destination"
	"github.com/unblee/cccp/source"
)

type worker struct {
	name  string
	src   source.Source
	dst   destination.Destination
	errCh chan error
	err   error
}

func newWorker(name string, src source.Source, dst destination.Destination, errCh chan error) *worker {
	return &worker{
		name:  name,
		src:   src,
		dst:   dst,
		errCh: errCh,
	}
}

func (wkr *worker) run(ctx context.Context) {
	defer wkr.dst.Close()
	defer wkr.src.Close()

	err := copyWithContext(ctx, wkr.dst, wkr.src)
	if err != nil {
		wkr.err = errors.Wrapf(err, "failed to copy src to dst '%s'", wkr.name)
		wkr.errCh <- wkr.err
		return
	}
}

// cf. http://ixday.github.io/post/golang-cancel-copy/
func copyWithContext(ctx context.Context, dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
	return err
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }
