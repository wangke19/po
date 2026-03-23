package exportcmd

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdExport(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <command>",
		Short: "Export Polarion data to CSV or JSON",
	}

	cmd.AddCommand(NewCmdWorkitems(f))
	cmd.AddCommand(NewCmdTestresults(f))
	return cmd
}
