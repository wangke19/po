// Package search provides the search command for querying Polarion workitems.
package search

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

// NewCmdSearch returns the 'search' command.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var wiType, status, author, jsonFields string
	var limit int

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search work items with optional Lucene query and filters",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			var parts []string
			if cmd.Flags().Changed("type") {
				parts = append(parts, "type:"+wiType)
			}
			if cmd.Flags().Changed("status") {
				parts = append(parts, "status:"+status)
			}
			if cmd.Flags().Changed("author") {
				parts = append(parts, "author:"+author)
			}
			if len(args) > 0 && args[0] != "" {
				parts = append(parts, args[0])
			}
			q := strings.Join(parts, " AND ")

			items, err := client.ListWorkItems(cmd.Context(), q, limit)
			if err != nil {
				return fmt.Errorf("search: %w", err)
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
				_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\t%s\n", item.ID, item.Type, item.Status, item.Title)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&wiType, "type", "", "Filter by work item type")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&author, "author", "", "Filter by author")
	cmd.Flags().IntVar(&limit, "limit", 50, "Max results")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
