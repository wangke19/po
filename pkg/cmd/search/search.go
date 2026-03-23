package search

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search work items with a Lucene query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			items, err := client.ListWorkItems(cmd.Context(), args[0], limit)
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
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, item := range items {
				fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\t%s\n", item.ID, item.Type, item.Status, item.Title)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 30, "Max results")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
