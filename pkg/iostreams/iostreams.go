// Package iostreams provides I/O stream abstractions for CLI interactions.
package iostreams

import (
	"io"
	"os"
	"strings"
)

// IOStreams represents the standard I/O streams.
type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer
}

// System returns the standard system I/O streams.
func System() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// IsTerminal returns whether the output is a terminal.
func (s *IOStreams) IsTerminal() bool {
	if f, ok := s.Out.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// Test returns I/O streams suitable for testing.
func Test() *IOStreams {
	return &IOStreams{
		In:     io.NopCloser(strings.NewReader("")),
		Out:    io.Discard,
		ErrOut: io.Discard,
	}
}
