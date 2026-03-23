package testcase

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdStepRemove(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "step-remove <case-id> <step-index>",
		Short: "Remove a test step by index",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			stepIndex, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid step index %q: %w", args[1], err)
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			steps, err := client.DeleteTestStep(cmd.Context(), args[0], stepIndex)
			if err != nil {
				return fmt.Errorf("remove step: %w", err)
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

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
