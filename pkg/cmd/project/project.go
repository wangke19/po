package project

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdProject returns the 'project' command group.
func NewCmdProject(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "Manage Polarion projects",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	return cmd
}
