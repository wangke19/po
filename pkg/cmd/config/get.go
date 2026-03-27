package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdGet returns the 'config get' command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <host> <key>",
		Short: "Get a configuration value for a host",
		Long: `Get a configuration value for a Polarion host.

Keys:
  project     Default project ID
  verify-ssl  TLS verification: true or false`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			hostname := config.NormalizeHostname(args[0])
			key := args[1]

			cfg, err := f.Config()
			if err != nil {
				return err
			}

			switch key {
			case "project":
				project, err := cfg.DefaultProject(hostname)
				if err != nil {
					return err
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, project)
			case "verify-ssl":
				_, _ = fmt.Fprintln(f.IOStreams.Out, cfg.VerifySSL(hostname))
			default:
				return fmt.Errorf("unknown key %q: valid keys are project, verify-ssl", key)
			}

			return nil
		},
	}
}
