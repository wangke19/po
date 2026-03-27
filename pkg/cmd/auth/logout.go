package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/zalando/go-keyring"
)

// NewCmdLogout returns the 'auth logout' command.
func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	var hostname string
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of a Polarion instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			if hostname == "" {
				hostname, err = cfg.DefaultHost()
				if err != nil {
					return err
				}
			}
			hostname = config.NormalizeHostname(hostname)
			_ = keyring.Delete("po", hostname)
			if err := cfg.RemoveHost(hostname); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(f.IOStreams.Out, "Logged out of %s\n", hostname)
			return nil
		},
	}
	cmd.Flags().StringVar(&hostname, "hostname", "", "Polarion instance hostname")
	return cmd
}
