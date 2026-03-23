package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdResult(f *cmdutil.Factory) *cobra.Command {
	var result, comment, jsonFields string
	cmd := &cobra.Command{
		Use:   "result <run-id> <case-id>",
		Short: "Record a test result",
		Args:  cobra.ExactArgs(2),
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

			if err := client.UpdateTestRunResult(cmd.Context(), args[0], args[1], polarion.TestResult{
				Result:  result,
				Comment: comment,
			}); err != nil {
				return fmt.Errorf("record test result: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				r := polarion.TestResult{Result: result, Comment: comment}
				out, err := jsonfields.FilterFields(r, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&result, "result", "", "Result: passed|failed|blocked (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Optional comment")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Print result as JSON (optional field list)")
	_ = cmd.MarkFlagRequired("result")
	return cmd
}
