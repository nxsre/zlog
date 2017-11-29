package main

import (
	"os"

	"github.com/soopsio/zlog/zlogbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
