package repos

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cmdutil "github.com/yhinoz/git-org/cmd/util"
)

func NewReposCmd(f cmdutil.Factory) *cobra.Command {
	c := &cobra.Command{
		Use:   "repos",
		Short: "Show the specified github organization repository",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := cmd.MarkFlagRequired("org"); err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			client, err := f.GithubClient()
			if err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}

			org, _ := cmd.Root().PersistentFlags().GetString("org")
			allRepos, err := cmdutil.GetAllRepository(client, org)
			if err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}

			for _, repo := range allRepos {
				fmt.Println(*repo.Name)
			}
		},
	}

	return c
}
