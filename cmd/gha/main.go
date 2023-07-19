package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var app = &cli.App{
	Name:    "gha",
	Usage:   "GitHub Analyzer",
	Version: "0.1.0",
	Commands: []*cli.Command{
		snapshotCommand,
		diffCommand,
		reportCommand,
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
