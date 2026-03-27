// Package importcmd provides import commands for workitems and test results.
package importcmd

import (
	"fmt"
	"io"
	"os"
)

func openReader(path string) (io.Reader, func(), error) {
	if path == "" || path == "-" {
		return os.Stdin, func() {}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open input file: %w", err)
	}
	return f, func() { _ = f.Close() }, nil
}
