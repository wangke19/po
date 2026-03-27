package workitem

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/browser"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

// NewCmdView returns the 'workitem view' command.
func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	var web bool
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View a work item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			item, err := client.GetWorkItem(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get work item %q: %w", args[0], err)
			}

			if web {
				return browser.Open(item.URL)
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

			_, _ = fmt.Fprintf(f.IOStreams.Out, "ID:          %s\nTitle:       %s\nType:        %s\nStatus:      %s\nAuthor:      %s\nDescription: %s\nURL:         %s\n",
				item.ID, item.Title, item.Type, item.Status, item.Author, item.Description, item.URL)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&web, "web", "w", false, "Open in browser")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON with specified fields")
	return cmd
}
