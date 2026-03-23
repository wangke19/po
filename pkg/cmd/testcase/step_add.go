package testcase

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdStepAdd(f *cmdutil.Factory) *cobra.Command {
	var action, expectedResult, jsonFields string

	cmd := &cobra.Command{
		Use:   "step-add <case-id>",
		Short: "Add a test step to a test case",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			steps, err := client.AddTestStep(cmd.Context(), args[0], polarion.TestStepInput{
				Action:         action,
				ExpectedResult: expectedResult,
			})
			if err != nil {
				return fmt.Errorf("add step: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(steps, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, s := range steps {
				fmt.Fprintf(f.IOStreams.Out, "%d\t%s\t%s\n", s.StepIndex, s.Action, s.ExpectedResult)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&action, "action", "", "Test step action (required)")
	cmd.Flags().StringVar(&expectedResult, "expected-result", "", "Expected result")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	_ = cmd.MarkFlagRequired("action")
	return cmd
}
