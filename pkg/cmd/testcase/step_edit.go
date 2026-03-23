package testcase

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdStepEdit(f *cmdutil.Factory) *cobra.Command {
	var action, expectedResult, jsonFields string

	cmd := &cobra.Command{
		Use:   "step-edit <case-id> <step-index>",
		Short: "Edit a test step's action or expected result",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("action") && !cmd.Flags().Changed("expected-result") {
				return fmt.Errorf("at least one of --action or --expected-result must be provided")
			}

			stepIndex, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid step index %q: %w", args[1], err)
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			steps, err := client.UpdateTestStep(cmd.Context(), args[0], stepIndex, polarion.TestStepInput{
				Action:         action,
				ExpectedResult: expectedResult,
			})
			if err != nil {
				return fmt.Errorf("edit step: %w", err)
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

	cmd.Flags().StringVar(&action, "action", "", "New action text")
	cmd.Flags().StringVar(&expectedResult, "expected-result", "", "New expected result text")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
