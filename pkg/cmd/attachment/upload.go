package attachment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdUpload(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "upload <work-item-id> <file>",
		Short: "Upload a file as an attachment to a work item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workItemID := args[0]
			filePath := args[1]

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer func() { _ = file.Close() }()

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			att, err := client.UploadAttachment(cmd.Context(), workItemID, filepath.Base(filePath), file)
			if err != nil {
				return fmt.Errorf("upload: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(att, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\t%d\n", att.ID, att.FileName, att.ContentType, att.Size)
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
