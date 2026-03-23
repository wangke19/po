package pocmd

import (
	"fmt"
	"os"

	"github.com/wangke19/po/internal/build"
	"github.com/wangke19/po/pkg/cmd/root"
	"github.com/wangke19/po/pkg/cmdutil"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func Main() exitCode {
	f := cmdutil.New(build.Version)
	cmd := root.NewCmdRoot(f, build.Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}
	return exitOK
}
