package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/urfave/cli/v3"
)

var snapshotCommand = &cli.Command{
	Name:      "snapshot",
	ArgsUsage: "<org>/<repo>",
	Usage:     "take a snapshot of a repository",
	Aliases:   []string{"s"},
	Action:    runSnapshot,
}

func runSnapshot(ctx *cli.Context) error {
	ref := ctx.Args().First()
	org, repo, ok := strings.Cut(ref, "/")
	if !ok {
		return fmt.Errorf("invalid ref: %s", ref)
	}

	client := github.NewClient()
	client.PageEvent = func(page int) {
		fmt.Printf("Fetching page %d...\n", page)
	}
	snapshot, n, err := client.Snapshot(ctx.Context, org, repo)
	if err != nil {
		return err
	}
	fmt.Println("Fetched", n, "issues and pull requests")

	path := fmt.Sprintf("%s_%s_%s.json", org, repo, time.Now().Format("20060102_150405"))
	if err := os.WriteFile(path, snapshot, 0644); err != nil {
		return err
	}
	fmt.Println("Saved snapshot to", path)

	return nil
}
