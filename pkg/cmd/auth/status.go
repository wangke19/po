package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/internal/config"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
	"github.com/zalando/go-keyring"
)

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
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

			type hostStatus struct {
				Host        string `json:"host"`
				Project     string `json:"project"`
				TokenStatus string `json:"tokenStatus"`
			}

			statuses := make([]hostStatus, 0, len(hosts))
			for _, host := range hosts {
				proj, _ := cfg.DefaultProject(host)
				ts := "token set"
				if os.Getenv("POLARION_TOKEN") != "" {
					ts = "env"
				} else {
					_, kerr := keyring.Get("po", host)
					if kerr != nil {
						ts = "no token"
					}
				}
				statuses = append(statuses, hostStatus{Host: host, Project: proj, TokenStatus: ts})
			}

			if cmd.Flags().Changed("json") {
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(statuses, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			for _, s := range statuses {
				fmt.Fprintf(f.IOStreams.Out, "%s\n  Project: %s\n  Token: %s\n", s.Host, s.Project, s.TokenStatus)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
