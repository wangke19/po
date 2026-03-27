package attachment

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdAttachment returns the 'attachment' command.
func NewCmdAttachment(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment <command>",
		Short: "Manage work item attachments",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdUpload(f))
	cmd.AddCommand(NewCmdDownload(f))
	return cmd
}
