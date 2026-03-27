// Package whoami provides the whoami command to display the current authenticated user.
package whoami

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
	"github.com/wangke19/po/pkg/jsonfields"
)

// NewCmdWhoami returns the 'whoami' command.
func NewCmdWhoami(f *cmdutil.Factory) *cobra.Command {
	var jsonFields string

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Display the current authenticated user and context",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			host, err := cfg.DefaultHost()
			if err != nil {
				return err
			}
			project, err := cfg.DefaultProject(host)
			if err != nil {
				return err
			}

			client, err := f.PolarionClient()
			if err != nil {
				return err
			}

			user, err := client.GetCurrentUser(cmd.Context())
			if err != nil {
				return err
			}

			projectURL := fmt.Sprintf("https://%s/polarion/#/project/%s", host, project)

			if cmd.Flags().Changed("json") {
				type output struct {
					Host       string `json:"host"`
					Project    string `json:"project"`
					ProjectURL string `json:"projectUrl"`
					ID         string `json:"id"`
					Name       string `json:"name"`
					Email      string `json:"email"`
				}
				fields := strings.Split(jsonFields, ",")
				if jsonFields == "" {
					fields = nil
				}
				out, err := jsonfields.FilterFields(output{
					Host:       host,
					Project:    project,
					ProjectURL: projectURL,
					ID:         user.ID,
					Name:       user.Name,
					Email:      user.Email,
				}, fields)
				if err != nil {
					return fmt.Errorf("filter fields: %w", err)
				}
				_, _ = fmt.Fprintln(f.IOStreams.Out, string(out))
				return nil
			}

			_, _ = fmt.Fprintf(f.IOStreams.Out, "Logged in to %s as %s (project: %s)\n", host, user.ID, project)
			if user.Name != "" {
				_, _ = fmt.Fprintf(f.IOStreams.Out, "Name:    %s\n", user.Name)
			}
			if user.Email != "" {
				_, _ = fmt.Fprintf(f.IOStreams.Out, "Email:   %s\n", user.Email)
			}
			_, _ = fmt.Fprintf(f.IOStreams.Out, "Project: %s\n", projectURL)
			return nil
		},
	}

	cmd.Flags().StringVar(&jsonFields, "json", "", "Output as JSON with specified fields (comma-separated)")
	return cmd
}
