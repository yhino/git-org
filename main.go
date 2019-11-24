package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/yhino/git-org/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		_ = godotenv.Load(os.Getenv("HOME") + "/.env")
	}

	defaultCmd := cmd.NewDefaultCmd()
	if err := defaultCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
