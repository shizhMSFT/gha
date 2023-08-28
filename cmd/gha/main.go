package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v3"
)

var app = &cli.Command{
	Name:    "gha",
	Usage:   "GitHub Analyzer",
	Version: "0.3.0",
	Commands: []*cli.Command{
		snapshotCommand,
		diffCommand,
		reportCommand,
		pullRequestReviewCommand,
		issueCommentCommand,
	},
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
