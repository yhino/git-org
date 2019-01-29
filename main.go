package main

import (
	"fmt"
	"os"

	"github.com/yhinoz/git-org/cmd"
)

func main() {
	defaultCmd := cmd.NewDefaultCmd()
	if err := defaultCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
