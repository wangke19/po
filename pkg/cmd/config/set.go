package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdSet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "set <host> <key> <value>",
		Short: "Set a configuration value for a host",
		Long: `Set a configuration value for a Polarion host.

Keys:
  project     Default project ID
  verify-ssl  TLS verification: true or false`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname := config.NormalizeHostname(args[0])
			key := args[1]
			value := args[2]

			cfg, err := f.Config()
			if err != nil {
				return err
			}

			switch key {
			case "project":
				verifySSL := cfg.VerifySSL(hostname)
				if err := cfg.SetHost(hostname, value, verifySSL); err != nil {
					return fmt.Errorf("set project: %w", err)
				}
			case "verify-ssl":
				v, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("verify-ssl must be true or false")
				}
				project, _ := cfg.DefaultProject(hostname)
				if err := cfg.SetHost(hostname, project, v); err != nil {
					return fmt.Errorf("set verify-ssl: %w", err)
				}
			default:
				return fmt.Errorf("unknown key %q: valid keys are project, verify-ssl", key)
			}

			fmt.Fprintf(f.IOStreams.Out, "Set %s.%s = %s\n", hostname, key, value)
			return nil
		},
	}
}
