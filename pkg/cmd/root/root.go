package root

import (
	"github.com/spf13/cobra"
	apiCmd "github.com/wangke19/po/pkg/cmd/api"
	authCmd "github.com/wangke19/po/pkg/cmd/auth"
	casecmd "github.com/wangke19/po/pkg/cmd/case"
	testrunCmd "github.com/wangke19/po/pkg/cmd/testrun"
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
	return cmd
}
