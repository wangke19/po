package link

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdLink returns the 'link' command.
func NewCmdLink(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link <command>",
		Short: "Manage work item links",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdAdd(f))
	cmd.AddCommand(NewCmdRemove(f))
	return cmd
}
