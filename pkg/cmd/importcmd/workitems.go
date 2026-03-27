package importcmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/importer"
	"github.com/wangke19/po/pkg/polarion"
)

// NewCmdWorkitems returns the 'import workitems' command.
func NewCmdWorkitems(f *cmdutil.Factory) *cobra.Command {
	var file, format string

	cmd := &cobra.Command{
		Use:   "workitems",
		Short: "Import work items from CSV or JSON",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			r, closer, err := openReader(file)
			if err != nil {
				return err
			}
			defer closer()

			items, err := parseWorkItems(r, format)
			if err != nil {
				return err
			}

			for _, item := range items {
				created, err := client.CreateWorkItem(cmd.Context(), item)
				if err != nil {
					return fmt.Errorf("create work item %q: %w", item.Title, err)
				}
				_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", created.ID, created.Type, created.Title)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Input file path (default: stdin)")
	cmd.Flags().StringVar(&format, "format", "csv", "Input format: csv or json")
	return cmd
}

func parseWorkItems(r io.Reader, format string) ([]polarion.WorkItemInput, error) {
	switch format {
	case "json":
		return importer.ReadWorkItemsJSON(r)
	default:
		return importer.ReadWorkItemsCSV(r)
	}
}
