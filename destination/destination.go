package destination

import (
	"io"
)

type Destination interface {
	io.WriteCloser
	Written() uint64
}
