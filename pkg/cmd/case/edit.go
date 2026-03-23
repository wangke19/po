package casecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdEdit(f *cmdutil.Factory) *cobra.Command {
	var title, wiType, desc, status, jsonFields string

	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a test case",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("title") && !cmd.Flags().Changed("type") &&
				!cmd.Flags().Changed("description") && !cmd.Flags().Changed("status") {
				return fmt.Errorf("at least one of --title, --type, --description, or --status must be provided")
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			input := polarion.WorkItemInput{}
			if cmd.Flags().Changed("title") {
				input.Title = title
			}
			if cmd.Flags().Changed("type") {
				input.Type = wiType
			}
			if cmd.Flags().Changed("status") {
				input.Status = status
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
	cmd.Flags().StringVar(&wiType, "type", "", "New work item type")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "New description")
	cmd.Flags().StringVar(&status, "status", "", "New status")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON (optional field list)")
	return cmd
}
