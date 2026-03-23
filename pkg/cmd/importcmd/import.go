package importcmd

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdImport(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <command>",
		Short: "Import data into Polarion from CSV or JSON",
	}

	cmd.AddCommand(NewCmdWorkitems(f))
	cmd.AddCommand(NewCmdTestresults(f))
	return cmd
}
