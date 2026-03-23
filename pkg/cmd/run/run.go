package run

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdRun(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <command>",
		Short: "Manage test run execution",
	}

	cmd.AddCommand(NewCmdStatus(f))
	cmd.AddCommand(NewCmdFinish(f))
	return cmd
}
