package config

import (
	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage po configuration",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdSet(f))
	cmd.AddCommand(NewCmdUnset(f))
	return cmd
}
