package exportcmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/export"
)

// NewCmdWorkitems returns the 'export workitems' command.
func NewCmdWorkitems(f *cmdutil.Factory) *cobra.Command {
	var wiType, query, format, output string
	var limit int

	cmd := &cobra.Command{
		Use:   "workitems",
		Short: "Export work items to CSV or JSON",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			q := fmt.Sprintf("type:%s", wiType)
			if query != "" {
				q += " AND " + query
			}

			items, err := client.ListWorkItems(cmd.Context(), q, limit)
			if err != nil {
				return fmt.Errorf("list work items: %w", err)
			}

			w, closer, err := openWriter(output, f.IOStreams.Out)
			if err != nil {
				return err
			}
			defer func() {
				if cerr := closer(); cerr != nil && err == nil {
					err = cerr
				}
			}()

			switch format {
			case "json":
				if err := export.WriteWorkItemsJSON(w, items); err != nil {
					return fmt.Errorf("export json: %w", err)
				}
			default:
				if err := export.WriteWorkItemsCSV(w, items); err != nil {
					return fmt.Errorf("export csv: %w", err)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&wiType, "type", "", "Work item type (required)")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Additional Lucene query")
	cmd.Flags().StringVar(&format, "format", "csv", "Output format: csv or json")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	cmd.Flags().IntVar(&limit, "limit", 100, "Max results")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func openWriter(path string, fallback io.Writer) (io.Writer, func() error, error) {
	if path == "" || path == "-" {
		return fallback, func() error { return nil }, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("create output file: %w", err)
	}
	return f, f.Close, nil
}
