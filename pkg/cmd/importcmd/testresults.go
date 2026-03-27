package importcmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/importer"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdTestresults(f *cmdutil.Factory) *cobra.Command {
	var file, format string

	cmd := &cobra.Command{
		Use:   "testresults <run-id>",
		Short: "Import test results into a test run from CSV or JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			r, closer, err := openReader(file)
			if err != nil {
				return err
			}
			defer closer()

			records, err := parseTestResults(r, format)
			if err != nil {
				return err
			}

			runID := args[0]
			var ok, failed int
			for _, rec := range records {
				err := client.UpdateTestRunResult(cmd.Context(), runID, rec.CaseID, polarion.TestResult{
					Result:  rec.Result,
					Comment: rec.Comment,
				})
				if err != nil {
					fmt.Fprintf(f.IOStreams.ErrOut, "warning: update %s: %v\n", rec.CaseID, err)
					failed++
					continue
				}
				ok++
			}
			_, _ = fmt.Fprintf(f.IOStreams.Out, "imported %d records (%d failed)\n", ok, failed)
			if failed > 0 {
				return fmt.Errorf("%d records failed to import", failed)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Input file path (default: stdin)")
	cmd.Flags().StringVar(&format, "format", "csv", "Input format: csv or json")
	return cmd
}

func parseTestResults(r io.Reader, format string) ([]polarion.TestRecord, error) {
	switch format {
	case "json":
		return importer.ReadTestResultsJSON(r)
	default:
		return importer.ReadTestResultsCSV(r)
	}
}
