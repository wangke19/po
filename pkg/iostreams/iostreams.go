package iostreams

import (
	"io"
	"os"
	"strings"
)

type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer
}

func System() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

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

func Test() *IOStreams {
	return &IOStreams{
		In:     io.NopCloser(strings.NewReader("")),
		Out:    io.Discard,
		ErrOut: io.Discard,
	}
}
