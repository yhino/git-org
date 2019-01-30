package grep

import (
	"bufio"
	"bytes"
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
)

func NewGrepCmd(f cmdutil.Factory) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "grep",
		Short: "run grep the specified github organization repository",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := cmd.MarkFlagRequired("org"); err != nil {
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

						tmpDir := filepath.Join(os.TempDir(), "git-orgrep", *repo.Name)
						currDir, _ := filepath.Abs(".")

						// clone
						branch, _ := cmd.Flags().GetString("branch")
						exec.Command("git", "clone", "-b", branch, "--quiet", *repo.SSHURL, tmpDir).Run()

						// grep
						var out []byte
						if _, err := os.Stat(tmpDir); err != nil {
							out = []byte("git-clone error")
						} else {
							os.Chdir(tmpDir)

							cmdArgs := append([]string{"grep"}, args...)
							out, err = exec.Command("git", cmdArgs...).CombinedOutput()
							if err != nil {
								var waitStatus syscall.WaitStatus
								if exitError, ok := err.(*exec.ExitError); ok {
									waitStatus = exitError.Sys().(syscall.WaitStatus)
								}
								switch waitStatus.ExitStatus() {
								case 1:
									out = []byte("notfound")
								default:
									out = []byte(fmt.Sprintf("error: %s, %s", out, err))
								}
							}

							os.Chdir(currDir)
						}

						m.Lock()

						// format and output
						r := bytes.NewReader(out)
						scanner := bufio.NewScanner(r)
						for scanner.Scan() {
							fmt.Printf("%s\t%s\n", *repo.FullName, scanner.Text())
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
	}

	cmds.Flags().StringP("branch", "b", "master", "run grep on the specified branch")

	return cmds
}
