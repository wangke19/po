package casecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

// NewCmdCreate returns the 'case create' command.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var title, desc, status, jsonFields string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new test case",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			item, err := client.CreateWorkItem(cmd.Context(), polarion.WorkItemInput{
				Title:       title,
				Type:        "testcase",
				Status:      status,
				Description: desc,
			})
			if err != nil {
				return fmt.Errorf("create test case: %w", err)
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
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Created test case %s\n%s\n", item.ID, item.URL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Test case title (required)")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "Description")
	cmd.Flags().StringVar(&status, "status", "draft", "Initial status (draft|approved)")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON (optional field list)")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}
