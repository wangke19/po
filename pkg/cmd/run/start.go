package run

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdStart(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "start <run-id>",
		Short: "Mark a test run as in progress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			run, err := client.UpdateTestRunStatus(cmd.Context(), args[0], "inprogress")
			if err != nil {
				return fmt.Errorf("run start: %w", err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(run, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\t%s\n", run.ID, run.Status, run.Title)
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
