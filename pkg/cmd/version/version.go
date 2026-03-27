// Package version provides the version command.
package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdVersion returns the 'version' command.
func NewCmdVersion(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the po version",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			_, _ = fmt.Fprintf(f.IOStreams.Out, "po version %s\n", f.AppVersion)
		},
	}
}
