package testrun

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdAttachDownload returns the 'testrun attach-download' command.
func NewCmdAttachDownload(f *cmdutil.Factory) *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "attach-download <run-id> <attachment-id>",
		Short: "Download an attachment from a test run",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			rc, err := client.DownloadTestRunAttachment(cmd.Context(), args[0], args[1])
			if err != nil {
				return fmt.Errorf("download: %w", err)
			}
			defer func() { _ = rc.Close() }()

			var dst io.Writer
			if output == "" || output == "-" {
				dst = f.IOStreams.Out
			} else {
				file, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
				defer func() { _ = file.Close() }()
				dst = file
			}

			if _, err := io.Copy(dst, rc); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	return cmd
}
