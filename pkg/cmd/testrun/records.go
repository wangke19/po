package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdRecords(f *cmdutil.Factory) *cobra.Command {
	var caseFilter, resultFilter, jsonFields string
	var notRun bool

	cmd := &cobra.Command{
		Use:   "records <run-id>",
		Short: "List test records for a test run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			records, err := client.GetTestRunRecords(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get test run records %q: %w", args[0], err)
			}

			if cmd.Flags().Changed("case") {
				filtered := records[:0]
				for _, r := range records {
					if r.CaseID == caseFilter {
						filtered = append(filtered, r)
					}
				}
				records = filtered
			}

			if cmd.Flags().Changed("result") {
				filtered := records[:0]
				for _, r := range records {
					if r.Result == resultFilter {
						filtered = append(filtered, r)
					}
				}
				records = filtered
			}

			if notRun {
				filtered := records[:0]
				for _, r := range records {
					if r.Result == "" {
						filtered = append(filtered, r)
					}
				}
				records = filtered
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(records, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, r := range records {
				fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", r.CaseID, r.Result, r.Comment)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&caseFilter, "case", "", "Filter by test case ID")
	cmd.Flags().StringVar(&resultFilter, "result", "", "Filter by result: passed, failed, blocked")
	cmd.Flags().BoolVar(&notRun, "not-run", false, "Show only records not yet executed")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
