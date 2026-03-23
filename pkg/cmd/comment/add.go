package comment

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var body, jsonFields string

	cmd := &cobra.Command{
		Use:   "add <work-item-id>",
		Short: "Add a comment to a work item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			c, err := client.AddComment(cmd.Context(), args[0], body)
			if err != nil {
				return fmt.Errorf("add comment: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(c, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\t%s\n", c.ID, c.Author, c.Created, c.Body)
			return nil
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "Comment text (required)")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	_ = cmd.MarkFlagRequired("body")
	return cmd
}
