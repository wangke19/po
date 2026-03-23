package root

import (
	"github.com/spf13/cobra"
	apiCmd "github.com/wangke19/po/pkg/cmd/api"
	authCmd "github.com/wangke19/po/pkg/cmd/auth"
	casecmd "github.com/wangke19/po/pkg/cmd/case"
	completionCmd "github.com/wangke19/po/pkg/cmd/completion"
	configCmd "github.com/wangke19/po/pkg/cmd/config"
	openCmd "github.com/wangke19/po/pkg/cmd/open"
	runCmd "github.com/wangke19/po/pkg/cmd/run"
	searchCmd "github.com/wangke19/po/pkg/cmd/search"
	testrunCmd "github.com/wangke19/po/pkg/cmd/testrun"
	versionCmd "github.com/wangke19/po/pkg/cmd/version"
	workitemCmd "github.com/wangke19/po/pkg/cmd/workitem"
	"github.com/wangke19/po/pkg/cmdutil"
)

func NewCmdRoot(f *cmdutil.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "po <command> <subcommand> [flags]",
		Short:         "Polarion CLI",
		Long:          "Work seamlessly with Polarion ALM from the command line.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(apiCmd.NewCmdApi(f))
	cmd.AddCommand(casecmd.NewCmdCase(f))
	cmd.AddCommand(testrunCmd.NewCmdTestrun(f))
	cmd.AddCommand(workitemCmd.NewCmdWorkitem(f))
	cmd.AddCommand(searchCmd.NewCmdSearch(f))
	cmd.AddCommand(configCmd.NewCmdConfig(f))
	cmd.AddCommand(versionCmd.NewCmdVersion(f))
	cmd.AddCommand(completionCmd.NewCmdCompletion(f))
	cmd.AddCommand(runCmd.NewCmdRun(f))
	cmd.AddCommand(openCmd.NewCmdOpen(f))
	return cmd
}
