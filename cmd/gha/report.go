package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/math"
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
		start = time.Now().UTC().AddDate(0, 0, int(-ago))
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
	if !start.IsZero() {
		fmt.Printf("- Start Date: `%s`\n", start.Format(time.DateTime))
	}
	if !end.IsZero() {
		fmt.Printf("- End Date: `%s`\n", end.Format(time.DateTime))
	}
	report := analysis.NewReport(start, end)
	for _, path := range ctx.Args().Slice() {
		fmt.Println()
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
	fmt.Println()
	fmt.Println("## Summary")
	printSummary(report.Abstract())
	return nil
}

func printSummary(summary *analysis.Summary) {
	printRepositorySummary(summary.RepositorySummary)

	fmt.Println()
	fmt.Println("### Contributors")
	for name, summary := range summary.Authors {
		fmt.Println()
		fmt.Println("####", name)
		printRepositorySummary(summary)
	}
}

func printRepositorySummary(summary *analysis.RepositorySummary) {
	issue := summary.Issue
	fmt.Println("Issues")
	fmt.Println("- Total:", issue.Total)
	fmt.Println("  - Open:", issue.Open)
	fmt.Println("  - Closed:", issue.Closed)
	fmt.Println("- Time to close:")
	fmt.Println("  - Min:", math.Min(issue.Durations))
	fmt.Println("  - Max:", math.Max(issue.Durations))
	fmt.Println("  - Mean:", math.Mean(issue.Durations))
	fmt.Println("  - Median:", math.Median(issue.Durations))
	fmt.Println("  - 90th Percentile:", math.Percentile(issue.Durations, 0.9))
	fmt.Println("  - 95th Percentile:", math.Percentile(issue.Durations, 0.95))
	fmt.Println("  - 99th Percentile:", math.Percentile(issue.Durations, 0.99))

	fmt.Println()

	pr := summary.PullRequest
	fmt.Println("Pull Requests")
	fmt.Println("- Total:", pr.Total)
	fmt.Println("  - Open:", pr.Open)
	fmt.Println("  - Closed:", pr.Closed)
	fmt.Println("  - Merged:", pr.Merged)
	fmt.Println("- Time to merge:")
	fmt.Println("  - Min:", math.Min(pr.Durations))
	fmt.Println("  - Max:", math.Max(pr.Durations))
	fmt.Println("  - Mean:", math.Mean(pr.Durations))
	fmt.Println("  - Median:", math.Median(pr.Durations))
	fmt.Println("  - 90th Percentile:", math.Percentile(pr.Durations, 0.9))
	fmt.Println("  - 95th Percentile:", math.Percentile(pr.Durations, 0.95))
	fmt.Println("  - 99th Percentile:", math.Percentile(pr.Durations, 0.99))
}
