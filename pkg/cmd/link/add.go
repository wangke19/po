// Package link provides commands for managing Polarion workitem links.
package link

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdAdd returns the 'link add' command.
func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var role string

	cmd := &cobra.Command{
		Use:   "add <work-item-id> <target-id>",
		Short: "Add a link between two work items",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			if err := client.AddLink(cmd.Context(), args[0], args[1], role); err != nil {
				return fmt.Errorf("add link: %w", err)
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Linked %s -> %s (role: %s)\n", args[0], args[1], role)
			return nil
		},
	}

	cmd.Flags().StringVar(&role, "role", "relates_to", "Link role (e.g. relates_to, blocks, duplicates)")
	return cmd
}
