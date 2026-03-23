package testrun

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/browser"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	var web bool
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View a test run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}
			run, err := client.GetTestRun(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get test run %q: %w", args[0], err)
			}

			if web {
				return browser.Open(run.URL)
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
			fmt.Fprintf(f.IOStreams.Out, "ID:       %s\nTitle:    %s\nStatus:   %s\nTemplate: %s\nURL:      %s\n",
				run.ID, run.Title, run.Status, run.Template, run.URL)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&web, "web", "w", false, "Open in browser")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON with specified fields")
	return cmd
}
