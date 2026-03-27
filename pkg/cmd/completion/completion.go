package completion

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdCompletion returns the 'completion' command.
func NewCmdCompletion(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "completion <shell>",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for po.

Supported shells: bash, zsh, fish, powershell

To load completions:

  Bash:
    source <(po completion bash)

  Zsh:
    po completion zsh > "${fpath[1]}/_po"

  Fish:
    po completion fish | source

  PowerShell:
    po completion powershell | Out-String | Invoke-Expression`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(f.IOStreams.Out)
			case "zsh":
				return root.GenZshCompletion(f.IOStreams.Out)
			case "fish":
				return root.GenFishCompletion(f.IOStreams.Out, true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(f.IOStreams.Out)
			default:
				return fmt.Errorf("unsupported shell %q: choose bash, zsh, fish, or powershell", args[0])
			}
		},
	}
}
