package exportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/export"
)

func NewCmdTestresults(f *cmdutil.Factory) *cobra.Command {
	var format, output string

	cmd := &cobra.Command{
		Use:   "testresults <run-id>",
		Short: "Export test run records to CSV or JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			records, err := client.GetTestRunRecords(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get test run records: %w", err)
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
				if err := export.WriteTestResultsJSON(w, records); err != nil {
					return fmt.Errorf("export json: %w", err)
				}
			default:
				if err := export.WriteTestResultsCSV(w, records); err != nil {
					return fmt.Errorf("export csv: %w", err)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "csv", "Output format: csv or json")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	return cmd
}
