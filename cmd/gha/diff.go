package main

import (
	"fmt"
	"os"

	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/urfave/cli/v3"
)

var diffCommand = &cli.Command{
	Usage:     "diff two snapshots",
	ArgsUsage: "<old> <new>",
	Aliases:   []string{"d"},
	Action:    runDiff,
}

func runDiff(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	oldPath := ctx.Args().Get(0)
	oldJSON, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}
	oldIssues, err := github.ParseIssues(oldJSON)
	if err != nil {
		return err
	}
	newPath := ctx.Args().Get(1)
	newJSON, err := os.ReadFile(newPath)
	if err != nil {
		return err
	}
	newIssues, err := github.ParseIssues(newJSON)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s: ", oldPath, newPath)

	diffs := github.DiffIssues(oldIssues, newIssues)
	if len(diffs) == 0 {
		fmt.Println("no change")
		return nil
	}

	fmt.Println(len(diffs), "changes:")
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
