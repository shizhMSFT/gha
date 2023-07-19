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
			Usage:    "only include snapshots that are at least `DAYS` old",
			OnlyOnce: true,
		},
		&cli.TimestampFlag{
			Name:     "start-date",
			Usage:    "only include snapshots that were created after `DATE`",
			Config:   cli.TimestampConfig{Layout: time.DateOnly},
			OnlyOnce: true,
		},
		&cli.TimestampFlag{
			Name:     "end-date",
			Usage:    "only include snapshots that were created before `DATE`",
			Config:   cli.TimestampConfig{Layout: time.DateOnly},
			OnlyOnce: true,
		},
	},
	Action: runReport,
}

func runReport(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return fmt.Errorf("invalid number of arguments")
	}

	// parse flags
	var start, end time.Time
	if ago := ctx.Int("ago"); ago > 0 {
		start = time.Now().UTC().AddDate(0, 0, -ago)
	}
	if date := ctx.Timestamp("start-date"); !date.IsZero() {
		start = *date
	}
	if date := ctx.Timestamp("end-date"); !date.IsZero() {
		end = *date
	}

	// generate report
	fmt.Println("GitHub Analysis Report")
	fmt.Println("======================")
	report := analysis.NewReport(start, end)
	for _, path := range ctx.Args().Slice() {
		fmt.Println("##", path)
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
	fmt.Println("## Summary")
	printSummary(report.Abstract())
	return nil
}

func printSummary(summary *analysis.Summary) {
	// TODO: print summary
}
