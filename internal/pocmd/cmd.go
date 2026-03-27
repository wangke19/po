package pocmd

import (
	"fmt"
	"os"

	"github.com/wangke19/po/internal/build"
	"github.com/wangke19/po/pkg/cmd/root"
	"github.com/wangke19/po/pkg/cmdutil"
)

// ExitCode represents the exit status of the command.
type ExitCode int

const (
	// ExitOK indicates successful execution.
	ExitOK ExitCode = 0
	// ExitError indicates an error occurred.
	ExitError ExitCode = 1
)

// Main executes the root command and returns the exit code.
func Main() ExitCode {
	f := cmdutil.New(build.Version)
	cmd := root.NewCmdRoot(f, build.Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitError
	}
	return ExitOK
}
