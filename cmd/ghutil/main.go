package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shizhMSFT/ghutil/pkg/github"
	"github.com/urfave/cli/v3"
)

var app = &cli.App{
	Name:  "ghutil",
	Usage: "GitHub utility",
	Commands: []*cli.Command{
		{
			Name:    "snapshot",
			Usage:   "create a snapshot",
			Aliases: []string{"s"},
			Commands: []*cli.Command{
				{
					Name:      "issue",
					ArgsUsage: "<org>/<repo>",
					Usage:     "create a snapshot of issues",
					Aliases:   []string{"i"},
					Action:    snapshotIssues,
				},
			},
		},
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

func snapshotIssues(ctx *cli.Context) error {
	ref := ctx.Args().First()
	org, repo, ok := strings.Cut(ref, "/")
	if !ok {
		return fmt.Errorf("invalid ref: %s", ref)
	}
	issues, err := github.OpenIssues(ctx.Context, org, repo)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s_%s_issues_%s.json", org, repo, time.Now().Format("20060102_150405"))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(issues)
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
