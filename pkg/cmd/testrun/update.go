package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

// NewCmdUpdate returns the 'testrun update' command.
func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var title, template, jsonFields string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a test run's title or template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("title") && !cmd.Flags().Changed("template") {
				return fmt.Errorf("at least one of --title or --template must be provided")
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			input := polarion.TestRunInput{}
			if cmd.Flags().Changed("title") {
				input.Title = title
			}
			if cmd.Flags().Changed("template") {
				input.Template = template
			}

			run, err := client.UpdateTestRun(cmd.Context(), args[0], input)
			if err != nil {
				return fmt.Errorf("update test run %q: %w", args[0], err)
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

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Updated %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "New title")
	cmd.Flags().StringVar(&template, "template", "", "New template ID")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
