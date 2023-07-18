package main

import (
	"fmt"
	"os"

	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/urfave/cli/v3"
)

var app = &cli.App{
	Name:  "gha",
	Usage: "GitHub Analyzer",
	Commands: []*cli.Command{
		snapshotCommand,
		{
			Name:    "diff",
			Usage:   "diff two snapshots",
			Aliases: []string{"d"},
			Commands: []*cli.Command{
				{
					Name:      "issue",
					ArgsUsage: "<old> <new>",
					Usage:     "diff two snapshots of issues",
					Aliases:   []string{"i"},
					Action:    diffIssues,
				},
			},
		},
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func diffIssues(ctx *cli.Context) error {
	oldIssues, err := os.ReadFile(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	newIssues, err := os.ReadFile(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	diffs, err := github.DiffIssues(oldIssues, newIssues)
	if err != nil {
		return err
	}

	if len(diffs) == 0 {
		fmt.Println("No changes")
		return nil
	}

	fmt.Println("Changes:")
	for _, diff := range diffs {
		fmt.Println()
		fmt.Printf("%s \033[0;33m#%d\033[0m\n", diff.Item.Title, diff.Item.Number)
		fmt.Println("Link:", diff.Item.HTMLURL)
		for _, change := range diff.Changes {
			fmt.Printf("+-- %s: \033[0;31m%s\033[0m -> \033[0;32m%s\033[0m\n", change.Field, change.Old, change.New)
		}
	}
	return nil
}
