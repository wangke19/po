package attachment

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdDownload(f *cmdutil.Factory) *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "download <work-item-id> <attachment-id>",
		Short: "Download an attachment from a work item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			rc, err := client.DownloadAttachment(cmd.Context(), args[0], args[1])
			if err != nil {
				return fmt.Errorf("download: %w", err)
			}
			defer rc.Close()

			var dst io.Writer
			if output == "" || output == "-" {
				dst = f.IOStreams.Out
			} else {
				file, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
				defer file.Close()
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
