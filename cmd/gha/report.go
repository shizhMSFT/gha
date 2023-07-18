package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/urfave/cli/v3"
)

var reportCommand = &cli.Command{
	Name:      "report",
	Usage:     "generate a report from snapshots",
	ArgsUsage: "<snapshot> [...]",
	Aliases:   []string{"r", "summarize"},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "ago",
			Aliases:  []string{"a"},
			Usage:    "only include snapshots that are at least `DAYS` old",
			OnlyOnce: true,
		},
	},
	Action: runReport,
}

func runReport(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return fmt.Errorf("invalid number of arguments")
	}

	var start, end time.Time
	if ago := ctx.Int("ago"); ago > 0 {
		end = time.Now().UTC()
		start = end.AddDate(0, 0, -ago)
		fmt.Println("Time frame:", start.Format(time.DateTime), "->", end.Format(time.DateTime))
	} else {
		fmt.Println("Time frame: all")
	}
	report := analysis.NewReport(start, end)
	for _, path := range ctx.Args().Slice() {
		fmt.Println(">>>", path)
		snapshotJSON, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot, err := github.ParseIssues(snapshotJSON)
		if err != nil {
			return err
		}
		printSummary(report.Summarize(path, snapshot))
	}
	fmt.Println("<<< Abstract")
	printSummary(report.Abstract())
	return nil
}

func printSummary(summary *analysis.Summary) {
	// TODO: print summary
}
