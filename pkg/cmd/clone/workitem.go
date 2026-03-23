package clone

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdWorkitem(f *cmdutil.Factory) *cobra.Command {
	var title, jsonFields string

	cmd := &cobra.Command{
		Use:   "workitem <id>",
		Short: "Clone a work item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			src, err := client.GetWorkItem(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get source work item: %w", err)
			}

			cloneTitle := title
			if cloneTitle == "" {
				cloneTitle = "Copy of " + src.Title
			}

			created, err := client.CreateWorkItem(cmd.Context(), polarion.WorkItemInput{
				Title:       cloneTitle,
				Type:        src.Type,
				Description: src.Description,
			})
			if err != nil {
				return fmt.Errorf("clone work item: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(created, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", created.ID, created.Type, created.Title)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Title for the cloned work item (default: \"Copy of <original title>\")")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
