/*
Package cccp provides a concurrent copy function with progress bars.
*/
package cccp

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/unblee/cccp/destination"
	"github.com/unblee/cccp/source"
)

// Run start concurrent copy.
func Run(ctx context.Context) error {
	argsLen := len(mngr.args)
	mngr.wkrs = make([]*worker, argsLen)
	mngr.pbs = make([]*progressbar, argsLen)
	mngr.errChs = make([]chan error, argsLen)
	wg := new(sync.WaitGroup)
	for i := 0; i < argsLen; i++ {
		wg.Add(1)
		arg := mngr.args[i]
		errCh := make(chan error, 1)
		mngr.wkrs[i] = newWorker(arg.name, arg.src, arg.dst, errCh)
		mngr.pbs[i] = newProgressbar(arg.name, arg.src.Size(), arg.dst, errCh)
		mngr.errChs[i] = errCh
	}

	limit := make(chan struct{}, mngr.concurrent)
	go func() {
		for i := 0; i < argsLen; i++ {
			limit <- struct{}{}
			go func(i int) {
				defer wg.Done()
				mngr.wkrs[i].run(ctx)
				close(mngr.errChs[i])
				<-limit
			}(i)
		}
	}()

	if mngr.disableProgressbar {
		wg.Wait()
	} else {
		mngr.printProgressbars(wg)
	}

	return mngr.composeErrors()
}

// SetFromSourceToDestination set the source and destination of the copy target.
// name is the name displayed in the progress bar.
func SetFromSourceToDestination(src source.Source, dst destination.Destination, name string) error {
	if name == "" {
		return errors.New("name is empty")
	}

	arg := &managerArg{
		src:  src,
		dst:  dst,
		name: name,
	}
	mngr.appendArg(arg)

	return nil
}

// SetFromURLToFile set the source URL and destination file path of the copy target.
// name is the name displayed in the progress bar.
// If name is empty, src URL base is set.
func SetFromURLToFile(src, dst, name string) error {
	switch {
	case src == "":
		return errors.New("src URL is empty")
	case dst == "":
		return errors.New("dst file path is empty")
	}

	srcURL, err := url.Parse(src)
	if err != nil {
		return errors.Wrapf(err, "invalid URL '%s'", srcURL)
	}

	if name == "" {
		name = path.Base(srcURL.Path)
	}

	s, err := source.NewURL(srcURL.String())
	if err != nil {
		return err
	}

	d, err := destination.NewFile(dst)
	if err != nil {
		return err
	}

	return SetFromSourceToDestination(s, d, name)
}

// SetFromFileToFile set the source file path and destination file path of the copy target.
// name is the name displayed in the progress bar.
// If name is empty, "src -> dst" is set.
func SetFromFileToFile(src, dst, name string) error {
	switch {
	case src == "":
		return errors.New("src URL is empty")
	case dst == "":
		return errors.New("dst file path is empty")
	case name == "":
		name = fmt.Sprintf("%s -> %s", src, dst)
	}

	s, err := source.NewFile(src)
	if err != nil {
		return err
	}

	d, err := destination.NewFile(dst)
	if err != nil {
		return err
	}

	return SetFromSourceToDestination(s, d, name)
}
