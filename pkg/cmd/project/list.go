package project

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all accessible Polarion projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			projects, err := client.ListProjects(cmd.Context())
			if err != nil {
				return err
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(projects, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, p := range projects {
				_, _ = fmt.Fprintf(f.IOStreams.Out, "%s\t%s\n", p.ID, p.Name)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
