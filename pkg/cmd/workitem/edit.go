package workitem

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdEdit(f *cmdutil.Factory) *cobra.Command {
	var title, desc, jsonFields string

	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a work item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			input := polarion.WorkItemInput{}
			if cmd.Flags().Changed("title") {
				input.Title = title
			}
			if cmd.Flags().Changed("description") {
				input.Description = desc
			}

			item, err := client.UpdateWorkItem(cmd.Context(), args[0], input)
			if err != nil {
				return fmt.Errorf("update work item %q: %w", args[0], err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(item, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "New title")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "New description")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON (optional field list)")
	return cmd
}
