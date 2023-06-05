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
				},
			},
			Action: snapshotIssues,
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
	path := fmt.Sprintf("%s_%s_issues_%s.json", org, repo, time.Now().Format("20060102"))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(issues)
}
