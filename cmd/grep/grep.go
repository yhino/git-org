package grep

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	cmdutil "github.com/yhinoz/git-org/cmd/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func NewGrepCmd(f cmdutil.Factory) *cobra.Command {
	tmpBaseDir := filepath.Join(os.TempDir(), "git-orgrep")

	cmds := &cobra.Command{
		Use:   "grep",
		Short: "Grep the specified github organization repository",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := cmd.MarkFlagRequired("org"); err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}

			if err := os.RemoveAll(tmpBaseDir); err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			service, err := f.RepositoryService()
			if err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}

			org, _ := cmd.Root().PersistentFlags().GetString("org")
			allRepos, err := service.FetchAllByOrg(org)
			if err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}

			thread := runtime.NumCPU()
			queue := make(chan *github.Repository, thread)

			var wg sync.WaitGroup
			var m sync.Mutex
			for i := 0; i <= thread; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					for {
						repo, ok := <-queue
						if !ok {
							return
						}

						tmpDir := filepath.Join(tmpBaseDir, *repo.Name)
						currDir, _ := filepath.Abs(".")

						out := func() []byte {
							// clone
							branch, _ := cmd.Flags().GetString("branch")
							if _, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
								URL:           *repo.SSHURL,
								ReferenceName: plumbing.NewBranchReferenceName(branch),
							}); err != nil {
								return []byte("clone error, " + err.Error())
							}

							m.Lock()

							os.Chdir(tmpDir)

							// grep
							var o []byte
							cmdArgs := append([]string{"grep"}, args...)
							o, err = exec.Command("git", cmdArgs...).CombinedOutput()
							if err != nil {
								var waitStatus syscall.WaitStatus
								if exitError, ok := err.(*exec.ExitError); ok {
									waitStatus = exitError.Sys().(syscall.WaitStatus)
								}
								switch waitStatus.ExitStatus() {
								case 1:
									o = []byte("notfound")
								default:
									o = []byte(fmt.Sprintf("error, %s, %s", o, err))
								}
							}

							os.Chdir(currDir)

							m.Unlock()

							return o
						}()

						m.Lock()

						// format and output
						results := ParseGrepResult(out)
						for _, result := range results {
							fmt.Printf("%s\t%s\n", *repo.FullName, result.TSV())
						}

						m.Unlock()
					}
				}()
			}
			for _, repo := range allRepos {
				queue <- repo
			}
			close(queue)
			wg.Wait()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := os.RemoveAll(tmpBaseDir); err != nil {
				cmd.Printf("ERROR: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmds.Flags().StringP("branch", "b", "master", "run grep on the specified branch")

	return cmds
}
