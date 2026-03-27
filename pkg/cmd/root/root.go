package root

import (
	"github.com/spf13/cobra"
	apiCmd "github.com/wangke19/po/pkg/cmd/api"
	attachmentCmd "github.com/wangke19/po/pkg/cmd/attachment"
	authCmd "github.com/wangke19/po/pkg/cmd/auth"
	casecmd "github.com/wangke19/po/pkg/cmd/case"
	cloneCmd "github.com/wangke19/po/pkg/cmd/clone"
	commentCmd "github.com/wangke19/po/pkg/cmd/comment"
	completionCmd "github.com/wangke19/po/pkg/cmd/completion"
	configCmd "github.com/wangke19/po/pkg/cmd/config"
	exportCmd "github.com/wangke19/po/pkg/cmd/exportcmd"
	importCmd "github.com/wangke19/po/pkg/cmd/importcmd"
	linkCmd "github.com/wangke19/po/pkg/cmd/link"
	openCmd "github.com/wangke19/po/pkg/cmd/open"
	projectCmd "github.com/wangke19/po/pkg/cmd/project"
	runCmd "github.com/wangke19/po/pkg/cmd/run"
	searchCmd "github.com/wangke19/po/pkg/cmd/search"
	testcaseCmd "github.com/wangke19/po/pkg/cmd/testcase"
	testrunCmd "github.com/wangke19/po/pkg/cmd/testrun"
	versionCmd "github.com/wangke19/po/pkg/cmd/version"
	whoamiCmd "github.com/wangke19/po/pkg/cmd/whoami"
	workitemCmd "github.com/wangke19/po/pkg/cmd/workitem"
	"github.com/wangke19/po/pkg/cmdutil"
)

// NewCmdRoot returns the root command for the po CLI.
func NewCmdRoot(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "po <command> <subcommand> [flags]",
		Short: "Polarion CLI",
		Long: `Work seamlessly with Polarion ALM from the command line.

Environment variables (take precedence over config file):
  POLARION_URL         Polarion server URL (e.g. https://polarion.example.com)
  POLARION_PROJECT     Default project ID (e.g. MYPROJECT)
  POLARION_TOKEN       Bearer token for authentication
  POLARION_VERIFY_SSL  Set to "false" to skip TLS certificate verification`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(apiCmd.NewCmdAPI(f))
	cmd.AddCommand(casecmd.NewCmdCase(f))
	cmd.AddCommand(testrunCmd.NewCmdTestrun(f))
	cmd.AddCommand(workitemCmd.NewCmdWorkitem(f))
	cmd.AddCommand(searchCmd.NewCmdSearch(f))
	cmd.AddCommand(configCmd.NewCmdConfig(f))
	cmd.AddCommand(versionCmd.NewCmdVersion(f))
	cmd.AddCommand(completionCmd.NewCmdCompletion(f))
	cmd.AddCommand(runCmd.NewCmdRun(f))
	cmd.AddCommand(openCmd.NewCmdOpen(f))
	cmd.AddCommand(testcaseCmd.NewCmdTestcase(f))
	cmd.AddCommand(attachmentCmd.NewCmdAttachment(f))
	cmd.AddCommand(linkCmd.NewCmdLink(f))
	cmd.AddCommand(commentCmd.NewCmdComment(f))
	cmd.AddCommand(cloneCmd.NewCmdClone(f))
	cmd.AddCommand(exportCmd.NewCmdExport(f))
	cmd.AddCommand(importCmd.NewCmdImport(f))
	cmd.AddCommand(projectCmd.NewCmdProject(f))
	cmd.AddCommand(whoamiCmd.NewCmdWhoami(f))
	return cmd
}
