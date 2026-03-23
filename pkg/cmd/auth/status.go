package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/zalando/go-keyring"
)

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			hosts := cfg.Hosts()
			if envURL := os.Getenv("POLARION_URL"); envURL != "" {
				hosts = []string{config.NormalizeHostname(envURL)}
			}
			if len(hosts) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "Not logged in to any Polarion instance.")
				return nil
			}
			for _, host := range hosts {
				proj, _ := cfg.DefaultProject(host)
				tokenStatus := "token set"
				if os.Getenv("POLARION_TOKEN") != "" {
					tokenStatus = "(from POLARION_TOKEN env)"
				} else {
					_, kerr := keyring.Get("po", host)
					if kerr != nil {
						tokenStatus = "no token"
					}
				}
				fmt.Fprintf(f.IOStreams.Out, "%s\n  Project: %s\n  Token: %s\n", host, proj, tokenStatus)
			}
			return nil
		},
	}
}
