package run

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdRun returns the 'run' command.
func NewCmdRun(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <command>",
		Short: "Manage test run execution",
	}

	cmd.AddCommand(NewCmdStatus(f))
	cmd.AddCommand(NewCmdStart(f))
	cmd.AddCommand(NewCmdPause(f))
	cmd.AddCommand(NewCmdFinish(f))
	return cmd
}
