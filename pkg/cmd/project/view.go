package project

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
		Short: "View a Polarion project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			p, err := client.GetProject(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get project %q: %w", args[0], err)
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(p, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			fmt.Fprintf(f.IOStreams.Out, "ID:          %s\nName:        %s\nDescription: %s\n",
				p.ID, p.Name, p.Description)
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
