package workitem

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdDelete returns the 'workitem delete' command.
func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a work item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("pass --confirm to delete work item %q", args[0])
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			if err := client.DeleteWorkItem(cmd.Context(), args[0]); err != nil {
				return fmt.Errorf("delete work item %q: %w", args[0], err)
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Deleted work item %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
	return cmd
}
