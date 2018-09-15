package cccp

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-tty"
	"github.com/unblee/cccp/destination"
	"github.com/unblee/cccp/source"
)

const (
	terminalCursorShow      = "\x1b[?25h"
	terminalCursorHide      = "\x1b[?25l"
	terminalClearScreenDown = "\x1b[0J"
)

var mngr = manager{
	mu:         new(sync.Mutex),
	concurrent: 1,
}

type managerArg struct {
	src  source.Source
	dst  destination.Destination
	name string
}

type manager struct {
	args                         []*managerArg
	wkrs                         []*worker
	pbs                          []*progressbar
	errChs                       []chan error
	mu                           *sync.Mutex
	concurrent                   int
	disableProgressbar           bool
	enableSequentialProgressbars bool
}

func (m *manager) appendArg(arg *managerArg) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mngr.args = append(mngr.args, arg)
}

func (m *manager) composeErrors() error {
	var errs error
	for _, wkr := range m.wkrs {
		if wkr.err != nil {
			errs = multierror.Append(errs, wkr.err)
		}
	}
	return errs
}

func (m *manager) printProgressbars(wg *sync.WaitGroup) {
	_tty, _ := tty.Open()
	defer _tty.Close()

	stdout := colorable.NewColorable(_tty.Output())
	fmt.Fprint(stdout, terminalCursorHide)
	defer fmt.Fprint(stdout, terminalCursorShow)

	doneCh := make(chan struct{}, 1)
	defer close(doneCh)

	go printProgressbarsLoop(wg, _tty, stdout, doneCh)

	wg.Wait()
	wg.Add(1)
	doneCh <- struct{}{}
	wg.Wait()

	// Displays that all execution has been completed.
	fmt.Fprint(stdout, terminalClearScreenDown)
	w, _, _ := _tty.Size()
	for _, pb := range mngr.pbs {
		fmt.Fprintln(stdout, pb.truncatedPrint(w))
	}
}

func printProgressbarsLoop(wg *sync.WaitGroup, _tty *tty.TTY, stdout io.Writer, doneCh chan struct{}) {
	var printedRows int
	var buf bytes.Buffer
	for {
		select {
		case <-doneCh:
			wg.Done()
			return
		default:
			printedRows = 0
			for i := 0; i < len(mngr.pbs); i++ {
				w, _, _ := _tty.Size()
				pb := mngr.pbs[i].truncatedPrint(w)
				if mngr.enableSequentialProgressbars && strings.HasPrefix(pb, waitFormatPrefix) {
					continue
				}
				buf.WriteString(pb + "\n")
				printedRows++
			}
			if printedRows >= 1 {
				buf.WriteString(fmt.Sprintf("\x1b[%dA", printedRows))
			}
			fmt.Fprint(stdout, buf.String())
			buf.Reset()
			time.Sleep(100 * time.Millisecond)
			fmt.Fprintf(stdout, "%s", terminalClearScreenDown)
		}
	}
}
