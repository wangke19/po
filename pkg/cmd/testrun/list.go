package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var status, query, jsonFields string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List test runs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			var parts []string
			if cmd.Flags().Changed("status") {
				parts = append(parts, "status:"+status)
			}
			if query != "" {
				parts = append(parts, query)
			}
			q := strings.Join(parts, " AND ")

			runs, err := client.ListTestRuns(cmd.Context(), q, limit)
			if err != nil {
				return fmt.Errorf("list test runs: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(runs, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, r := range runs {
				fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", r.ID, r.Status, r.Title)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Lucene query")
	cmd.Flags().IntVar(&limit, "limit", 30, "Max results")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
