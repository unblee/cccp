package cccp

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/unblee/cccp/destination"
)

const (
	terminalColorRed   = "\x1b[31m"
	terminalColorGreen = "\x1b[32m"
	terminalColorReset = "\x1b[m"
)

type progressbar struct {
	name        string
	size        uint64
	dst         destination.Destination
	spinCounter int
	done        bool
	errCh       chan error
	err         error
}

type progressbarFormat [3]string

func newProgressbar(name string, size uint64, dst destination.Destination, errCh chan error) *progressbar {
	return &progressbar{
		name:  name,
		size:  size,
		dst:   dst,
		errCh: errCh,
	}
}

func (pb *progressbar) truncatedPrint(termWidth int) string {
	format := pb.print()
	formatLen := len(format[0]) + len(format[1]) + len(format[2]) + 3 // +3 is separator(space) length

	// full format
	if formatLen <= termWidth {
		spaces := func() string {
			sp := ""
			for i := 0; i < termWidth-formatLen; i++ {
				sp += " "
			}
			return sp
		}()
		return fmt.Sprintf(
			"%s  %s %s%s",
			format[0],
			format[1],
			spaces,
			format[2],
		)
	}

	// short format
	formatLen = formatLen - (len(format[2]) + 2) // +2 is separator(space) length
	if formatLen <= termWidth {
		return fmt.Sprintf(
			"%s  %s",
			format[0],
			format[1],
		)
	}

	// truncated short format
	if overLen := formatLen - termWidth; overLen <= len(format[1])+2 { // +2 is a length of suffix that indicates truncation
		return fmt.Sprintf(
			"%s  %s%s",
			format[0],
			format[1][0:len(format[1])-overLen-1],
			"..",
		)
	}

	// shortest format
	return format[0]
}

func (pb *progressbar) print() progressbarFormat {
	select {
	case err, notClose := <-pb.errCh:
		if err != nil {
			pb.err = err
		}
		if !notClose {
			pb.done = true
		}
	default:
	}

	written := pb.dst.Written()
	switch {
	case pb.err != nil:
		return pb.errFormat()
	case pb.done:
		return pb.doneFormat()
	case written == 0:
		return pb.waitFormat()
	}

	return pb.downloadingFormat(written)
}

const waitFormatPrefix = "⚠ waiting..."

func (pb *progressbar) waitFormat() progressbarFormat {
	return progressbarFormat{
		waitFormatPrefix,
		pb.name,
		"",
	}
}

const errFormatPrefix = terminalColorRed + "❌ error: " + terminalColorReset

func (pb *progressbar) errFormat() progressbarFormat {
	return progressbarFormat{
		errFormatPrefix,
		pb.name,
		pb.err.Error(),
	}
}

const doneFormatPrefix = terminalColorGreen + "✅ complete!" + terminalColorReset

func (pb *progressbar) doneFormat() progressbarFormat {
	return progressbarFormat{
		doneFormatPrefix,
		pb.name,
		"",
	}
}

var spinSymbols = []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")

func (pb *progressbar) downloadingFormat(written uint64) progressbarFormat {
	symbol := string(spinSymbols[pb.spinCounter])
	pb.spinCounter++
	if pb.spinCounter >= len(spinSymbols) {
		pb.spinCounter = 0
	}
	return progressbarFormat{
		symbol + " downloading...",
		pb.name,
		fmt.Sprintf(
			"%s %3d%% %s",
			pb.progress(written),
			pb.parcent(written),
			pb.bar(written),
		),
	}
}

func (pb *progressbar) progress(written uint64) string {
	return fmt.Sprintf("%6s/%6s", humanize.Bytes(written), humanize.Bytes(pb.size))
}

func (pb *progressbar) bar(written uint64) string {
	var bar string
	bar += "["
	for i := 0; i < 10; i++ {
		if i < pb.parcent(written)/10 {
			bar += "#"
		} else {
			bar += "-"
		}
	}
	bar += "]"
	return bar
}

func (pb *progressbar) parcent(written uint64) int {
	parsent := float64(written) / float64(pb.size) * 100.0
	return int(parsent)
}
