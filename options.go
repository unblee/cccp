package cccp

import (
	"runtime"
)

// Option is a setting to change the behavior of manager.
type Option func(*manager)

// SetOptions set options to manager.
func SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt(&mngr)
	}
}

// WithConcurrent set the number of concurrent execution.
func WithConcurrent(n int) Option {
	return func(m *manager) {
		if n >= 1 {
			m.concurrent = n
		}
	}
}

// WithConcurrentNumCPU set the number of concurrent execution to the number of CPU cores.
func WithConcurrentNumCPU() Option {
	return func(m *manager) {
		m.concurrent = runtime.NumCPU()
	}
}

// WithDisableProgressbars set not to display the progress bars.
func WithDisableProgressbars() Option {
	return func(m *manager) {
		m.disableProgressbar = true
	}
}

// WithEnableSequentialProgressbars progress bars are displayed sequentially.
func WithEnableSequentialProgressbars() Option {
	return func(m *manager) {
		m.enableSequentialProgressbars = true
	}
}
