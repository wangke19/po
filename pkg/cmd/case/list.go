package casecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var query string
	var limit int
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List test cases",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			q := "type:testcase"
			if query != "" {
				q += " AND " + query
			}

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
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, item := range items {
				fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", item.ID, item.Status, item.Title)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Lucene query to filter results")
	cmd.Flags().IntVar(&limit, "limit", 30, "Maximum number of results")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
