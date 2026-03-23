package casecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View a test case",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			item, err := client.GetWorkItem(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get test case %q: %w", args[0], err)
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
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			fmt.Fprintf(f.IOStreams.Out, "ID:     %s\nTitle:  %s\nType:   %s\nStatus: %s\nURL:    %s\n",
				item.ID, item.Title, item.Type, item.Status, item.URL)
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output JSON with specified fields")
	return cmd
}
