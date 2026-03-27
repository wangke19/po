package casecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

// NewCmdTransition returns the 'case transition' command.
func NewCmdTransition(f *cmdutil.Factory) *cobra.Command {
	var to, jsonFields string

	cmd := &cobra.Command{
		Use:   "transition <id>",
		Short: "Transition a test case to a new status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			input := polarion.WorkItemInput{}
			if cmd.Flags().Changed("to") {
				input.Status = to
			}

			item, err := client.UpdateWorkItem(cmd.Context(), args[0], input)
			if err != nil {
				return fmt.Errorf("transition test case %q: %w", args[0], err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(item, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "%s transitioned to %s\n", args[0], item.Status)
			return nil
		},
	}

	cmd.Flags().StringVar(&to, "to", "", "Target status (required)")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON with specified fields (comma-separated)")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}
