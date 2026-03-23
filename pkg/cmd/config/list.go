package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured Polarion hosts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			hosts := cfg.Hosts()
			if len(hosts) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No hosts configured. Run: po auth login")
				return nil
			}

			for _, h := range hosts {
				project, _ := cfg.DefaultProject(h)
				verifySSL := cfg.VerifySSL(h)
				fmt.Fprintf(f.IOStreams.Out, "%s\tproject=%s\tverify-ssl=%v\n", h, project, verifySSL)
			}
			return nil
		},
	}
}
