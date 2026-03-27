package auth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/zalando/go-keyring"
)

// NewCmdToken returns the 'auth token' command.
func NewCmdToken(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Print the stored authentication token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if t := os.Getenv("POLARION_TOKEN"); t != "" {
				_, _ = fmt.Fprintln(f.IOStreams.Out, t)
				return nil
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			host, err := cfg.DefaultHost()
			if err != nil {
				return err
			}
			token, err := keyring.Get("po", host)
			if err != nil {
				return fmt.Errorf("no token stored for %s", host)
			}
			_, _ = fmt.Fprintln(f.IOStreams.Out, token)
			return nil
		},
	}
}
