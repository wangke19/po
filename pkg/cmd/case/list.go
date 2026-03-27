package casecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

// NewCmdList returns the 'case list' command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var status, author, query, jsonFields string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List test cases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			parts := []string{"type:testcase"}
			if cmd.Flags().Changed("status") {
				parts = append(parts, "status:"+status)
			}
			if cmd.Flags().Changed("author") {
				parts = append(parts, "author:"+author)
			}
			if query != "" {
				parts = append(parts, query)
			}
			q := strings.Join(parts, " AND ")

			items, err := client.ListWorkItems(cmd.Context(), q, limit)
			if err != nil {
				return fmt.Errorf("list test cases: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(items, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, item := range items {
				_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", item.ID, item.Status, item.Title)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&author, "author", "", "Filter by author")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Lucene query to filter results")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of results")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
