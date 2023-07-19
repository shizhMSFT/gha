package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/math"
	"github.com/shizhMSFT/gha/pkg/sort"
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
		&cli.BoolFlag{
			Name:     "contributors",
			Usage:    "include contributors from the report",
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
	if date := ctx.Value("start-date").(time.Time); !date.IsZero() {
		start = date
	}
	if date := ctx.Value("end-date").(time.Time); !date.IsZero() {
		end = date
	}
	includeContributors := ctx.Bool("contributors")

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
		printSummary(report.Summarize(path, snapshot), includeContributors)
	}
	fmt.Println()
	fmt.Println("## Summary")
	printSummary(report.Abstract(), includeContributors)
	return nil
}

func printSummary(summary *analysis.Summary, includeContributors bool) {
	printRepositorySummary(summary.RepositorySummary)

	if includeContributors {
		fmt.Println()
		fmt.Println("### Contributors")
		for name, summary := range summary.Authors {
			fmt.Println()
			fmt.Println("####", name)
			printRepositorySummary(summary)
		}
	}
}

func printRepositorySummary(summary *analysis.RepositorySummary) {
	issue := summary.Issue
	fmt.Println("Issues")
	fmt.Println("- Total:", issue.Total)
	fmt.Println("  - Open:", issue.Open)
	fmt.Println("  - Closed:", issue.Closed)
	if len(issue.Durations) > 0 {
		sort.Sort(issue.Durations)
		fmt.Println("- Time to close:")
		fmt.Println("  - Min:", formatDuration(math.Min(issue.Durations)))
		fmt.Println("  - Max:", formatDuration(math.Max(issue.Durations)))
		fmt.Println("  - Mean:", formatDuration(math.Mean(issue.Durations)))
		fmt.Println("  - Median:", formatDuration(math.Median(issue.Durations)))
		fmt.Println("  - 90th percentile:", formatDuration(math.Percentile(issue.Durations, 0.9)))
		fmt.Println("  - 95th percentile:", formatDuration(math.Percentile(issue.Durations, 0.95)))
		fmt.Println("  - 99th percentile:", formatDuration(math.Percentile(issue.Durations, 0.99)))
	}

	fmt.Println()

	pr := summary.PullRequest
	fmt.Println("Pull Requests")
	fmt.Println("- Total:", pr.Total)
	fmt.Println("  - Open:", pr.Open)
	fmt.Println("  - Closed:", pr.Closed)
	fmt.Println("  - Merged:", pr.Merged)
	if len(pr.Durations) > 0 {
		sort.Sort(pr.Durations)
		fmt.Println("- Time to merge:")
		fmt.Println("  - Min:", formatDuration(math.Min(pr.Durations)))
		fmt.Println("  - Max:", formatDuration(math.Max(pr.Durations)))
		fmt.Println("  - Mean:", formatDuration(math.Mean(pr.Durations)))
		fmt.Println("  - Median:", formatDuration(math.Median(pr.Durations)))
		fmt.Println("  - 90th percentile:", formatDuration(math.Percentile(pr.Durations, 0.9)))
		fmt.Println("  - 95th percentile:", formatDuration(math.Percentile(pr.Durations, 0.95)))
		fmt.Println("  - 99th percentile:", formatDuration(math.Percentile(pr.Durations, 0.99)))
	}
}

func formatDuration(d time.Duration) string {
	seconds := d / time.Second
	minutes, seconds := seconds/60, seconds%60
	hours, minutes := minutes/60, minutes%60
	days, hours := hours/24, hours%24
	months, days := days/30, days%30
	years, months := months/12, months%12
	if years > 0 {
		if months > 0 {
			return fmt.Sprintf("%dy %dm", years, months)
		}
		return fmt.Sprintf("%dy", years)
	}
	if months > 0 {
		if days > 0 {
			return fmt.Sprintf("%dmo %dd", months, days)
		}
		return fmt.Sprintf("%dmo", months)
	}
	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", seconds)
}
