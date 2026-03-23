package workitem

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var wiType, title, desc, jsonFields string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a work item",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			item, err := client.CreateWorkItem(cmd.Context(), polarion.WorkItemInput{
				Type:        wiType,
				Title:       title,
				Description: desc,
			})
			if err != nil {
				return fmt.Errorf("create work item: %w", err)
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

			fmt.Fprintf(f.IOStreams.Out, "Created %s %s\n%s\n", wiType, item.ID, item.URL)
			return nil
		},
	}

	cmd.Flags().StringVar(&wiType, "type", "", "Work item type (required)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Title (required)")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "Description")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON (optional field list)")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}
