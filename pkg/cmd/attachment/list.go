package attachment

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "list <work-item-id>",
		Short: "List attachments for a work item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			attachments, err := client.ListAttachments(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("list attachments: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(attachments, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, a := range attachments {
				_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\t%d\n", a.ID, a.FileName, a.ContentType, a.Size)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
