package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/wangke19/po/pkg/polarion"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var title, template, jsonFields string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new test run",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}
			run, err := client.CreateTestRun(cmd.Context(), polarion.TestRunInput{
				Title:    title,
				Template: template,
			})
			if err != nil {
				return fmt.Errorf("create test run: %w", err)
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
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}
			fmt.Fprintf(f.IOStreams.Out, "Created test run %s\n%s\n", run.ID, run.URL)
			return nil
		},
	}
	cmd.Flags().StringVarP(&title, "title", "t", "", "Test run title (required)")
	cmd.Flags().StringVar(&template, "template", "", "Template ID")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON (optional field list)")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}
