// Package open provides the open command to open Polarion resources in a browser.
package open

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/browser"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdOpen returns the 'open' command.
func NewCmdOpen(f *cmdutil.Factory) *cobra.Command {
	var resourceType string

	cmd := &cobra.Command{
		Use:   "open <id>",
		Short: "Open a Polarion resource in the browser",
		Long: `Open a Polarion work item, test run, or test case in the default browser.

Resource type is inferred from the ID prefix when --type is not specified:
  workitem  e.g. MYPROJ-123
  testrun   e.g. MyProject-run-001
  case      same as workitem (test cases are work items)

Use --type to override.`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"workitem", "testrun", "case"},
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			host, err := cfg.DefaultHost()
			if err != nil {
				return err
			}

			project, err := cfg.DefaultProject(host)
			if err != nil {
				return err
			}

			id := args[0]
			url := buildURL(host, project, id, resourceType)

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Opening %s\n", url)
			if err := browser.Open(url); err != nil {
				return fmt.Errorf("open browser: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&resourceType, "type", "workitem", "Resource type: workitem, testrun, case")
	return cmd
}

func buildURL(host, project, id, resourceType string) string {
	switch resourceType {
	case "testrun":
		return fmt.Sprintf("https://%s/polarion/#/project/%s/testrun?id=%s", host, project, id)
	default:
		return fmt.Sprintf("https://%s/polarion/#/project/%s/workitem?id=%s", host, project, id)
	}
}
