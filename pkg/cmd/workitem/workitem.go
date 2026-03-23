package workitem

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdWorkitem(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workitem <command>",
		Short: "Manage Polarion work items",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdEdit(f))
	cmd.AddCommand(NewCmdTransition(f))
	return cmd
}
