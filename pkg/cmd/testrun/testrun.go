package testrun

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdTestrun(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testrun <command>",
		Short: "Manage Polarion test runs",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdResult(f))
	cmd.AddCommand(NewCmdAddRecord(f))
	cmd.AddCommand(NewCmdRecords(f))
	return cmd
}
