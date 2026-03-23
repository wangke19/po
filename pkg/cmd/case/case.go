package casecmd

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdCase(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "case <command>",
		Short: "Manage Polarion test cases",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdEdit(f))
	cmd.AddCommand(NewCmdTransition(f))
	return cmd
}
