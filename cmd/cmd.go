package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhino/git-org/cmd/grep"
	"github.com/yhino/git-org/cmd/repos"
	cmdutil "github.com/yhino/git-org/cmd/util"
	"github.com/yhino/git-org/cmd/version"
)

func NewDefaultCmd() *cobra.Command {
	f := cmdutil.NewFactory()

	cmds := &cobra.Command{
		Use:   "git-org",
		Short: "A Git subcommand to do github orgaization",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmds.PersistentFlags().StringP("org", "o", "", "organization name")

	cmds.AddCommand(repos.NewReposCmd(f))
	cmds.AddCommand(grep.NewGrepCmd(f))
	cmds.AddCommand(version.NewVersionCmd())

	return cmds
}
