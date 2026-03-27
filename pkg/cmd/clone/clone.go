// Package clone provides commands for cloning Polarion resources.
package clone

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdClone returns the 'clone' command.
func NewCmdClone(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <command>",
		Short: "Clone a Polarion resource within the same project",
	}

	cmd.AddCommand(NewCmdWorkitem(f))
	cmd.AddCommand(NewCmdTestrun(f))
	return cmd
}
