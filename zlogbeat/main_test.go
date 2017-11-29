package main

// This file is mandatory as otherwise the zlogbeat.test binary is not generated correctly.

import (
	"flag"
	"testing"

	"github.com/soopsio/zlog/zlogbeat/cmd"
)

var systemTest *bool

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
	flag.String("c", "beat.yml", "Configuration file, relative to path.config")
	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup("systemTest"))
	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup("test.coverprofile"))
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {

	if *systemTest {
		main()
	}
}
