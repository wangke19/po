package run

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "status <run-id>",
		Short: "Show test run progress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			runID := args[0]
			records, err := client.GetTestRunRecords(cmd.Context(), runID)
			if err != nil {
				return fmt.Errorf("run status: %w", err)
			}

			progress := computeProgress(records)

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(progress, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Total:   %d\n", progress.Total)
			_, _ = fmt.Fprintf(f.IOStreams.Out, "Passed:  %d\n", progress.Passed)
			_, _ = fmt.Fprintf(f.IOStreams.Out, "Failed:  %d\n", progress.Failed)
			_, _ = fmt.Fprintf(f.IOStreams.Out, "Blocked: %d\n", progress.Blocked)
			_, _ = fmt.Fprintf(f.IOStreams.Out, "Not run: %d\n", progress.NotRun)
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}

func computeProgress(records []polarion.TestRecord) polarion.TestRunProgress {
	var p polarion.TestRunProgress
	p.Total = len(records)
	for _, r := range records {
		switch r.Result {
		case "passed":
			p.Passed++
		case "failed":
			p.Failed++
		case "blocked":
			p.Blocked++
		default:
			p.NotRun++
		}
	}
	return p
}
