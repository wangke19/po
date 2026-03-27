package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdUnset(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "unset <host>",
		Short: "Remove a host from configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname := config.NormalizeHostname(args[0])

			cfg, err := f.Config()
			if err != nil {
				return err
			}

			if err := cfg.RemoveHost(hostname); err != nil {
				return fmt.Errorf("unset host: %w", err)
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Removed %s from configuration\n", hostname)
			return nil
		},
	}
}
