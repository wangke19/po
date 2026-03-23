package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdAddRecord(f *cmdutil.Factory) *cobra.Command {
	var caseID, result, comment, jsonFields string

	cmd := &cobra.Command{
		Use:   "add-record <run-id>",
		Short: "Add a test record to a test run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch result {
			case "passed", "failed", "blocked":
			default:
				return fmt.Errorf("--result must be passed, failed, or blocked")
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			if err := client.AddTestRecord(cmd.Context(), args[0], caseID, polarion.TestResult{
				Result:  result,
				Comment: comment,
			}); err != nil {
				return err
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				r := polarion.TestRecord{CaseID: caseID, Result: result, Comment: comment}
				out, err := jsonfields.FilterFields(r, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			fmt.Fprintf(f.IOStreams.Out, "Added record for %s: %s\n", caseID, result)
			return nil
		},
	}

	cmd.Flags().StringVar(&caseID, "case", "", "Test case ID (required)")
	cmd.Flags().StringVar(&result, "result", "", "Result: passed|failed|blocked (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Optional comment")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON with specified fields (comma-separated)")
	_ = cmd.MarkFlagRequired("case")
	_ = cmd.MarkFlagRequired("result")
	return cmd
}
