package testcase

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdTestcase(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testcase <command>",
		Short: "Manage test case steps",
	}

	cmd.AddCommand(NewCmdSteps(f))
	cmd.AddCommand(NewCmdStepAdd(f))
	cmd.AddCommand(NewCmdStepRemove(f))
	cmd.AddCommand(NewCmdStepEdit(f))
	return cmd
}
