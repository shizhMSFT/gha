package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
		return errors.New("no snapshot files specified")
	}

	// parse flags
	var timeFrame analysis.TimeFrame
	if ago := ctx.Int("ago"); ago > 0 {
		timeFrame.Start = time.Now().UTC().AddDate(0, 0, int(-ago))
	}
	if date := ctx.Value("start-date").(time.Time); !date.IsZero() {
		timeFrame.Start = date
	}
	if date := ctx.Value("end-date").(time.Time); !date.IsZero() {
		timeFrame.End = date
	}
	includeContributors := ctx.Bool("contributors")

	// generate report
	fmt.Println("GitHub Analysis Report")
	fmt.Println("======================")
	printTimeFrame(timeFrame)
	report := analysis.NewReport(timeFrame)
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
	if ctx.NArg() > 1 {
		fmt.Println()
		fmt.Println("## Overall")
		printSummary(report.Abstract(), includeContributors)
	}
	return nil
}

func printTimeFrame(timeFrame analysis.TimeFrame) {
	if !timeFrame.Start.IsZero() {
		fmt.Printf("- Start Date: `%s`\n", timeFrame.Start.Format(time.DateTime))
	}
	if !timeFrame.End.IsZero() {
		fmt.Printf("- End Date: `%s`\n", timeFrame.End.Format(time.DateTime))
	}
}

func printSummary(summary *analysis.Summary, includeContributors bool) {
	printRepositorySummary(summary.RepositorySummary)

	if includeContributors {
		fmt.Println()
		fmt.Println("### Contributors")
		printContributors(summary.Authors)
	}
}

func printRepositorySummary(summary *analysis.RepositorySummary) {
	fmt.Println()
	fmt.Println("Issues")
	printIssueSummary(summary.Issue)

	fmt.Println()
	fmt.Println("Pull Requests")
	printPullRequestSummary(summary.PullRequest)
}

func printIssueSummary(summary *analysis.IssueSummary) {
	fmt.Println("- Total:", summary.Total)
	fmt.Println("  - Open:", summary.Open)
	fmt.Println("  - Closed:", summary.Closed)
	if len(summary.Durations) > 0 {
		sort.Sort(summary.Durations)
		fmt.Println("- Time to close:")
		fmt.Println("  - Min:", formatDuration(math.Min(summary.Durations)))
		fmt.Println("  - Max:", formatDuration(math.Max(summary.Durations)))
		fmt.Println("  - Mean:", formatDuration(math.Mean(summary.Durations)))
		fmt.Println("  - Median:", formatDuration(math.Median(summary.Durations)))
		fmt.Println("  - 90th percentile:", formatDuration(math.Percentile(summary.Durations, 0.9)))
		fmt.Println("  - 95th percentile:", formatDuration(math.Percentile(summary.Durations, 0.95)))
		fmt.Println("  - 99th percentile:", formatDuration(math.Percentile(summary.Durations, 0.99)))
	}
}

func printPullRequestSummary(summary *analysis.PullRequestSummary) {
	fmt.Println("- Total:", summary.Total)
	fmt.Println("  - Open:", summary.Open)
	fmt.Println("  - Closed:", summary.Closed)
	fmt.Println("  - Merged:", summary.Merged)
	if len(summary.Durations) > 0 {
		sort.Sort(summary.Durations)
		fmt.Println("- Time to merge:")
		fmt.Println("  - Min:", formatDuration(math.Min(summary.Durations)))
		fmt.Println("  - Max:", formatDuration(math.Max(summary.Durations)))
		fmt.Println("  - Mean:", formatDuration(math.Mean(summary.Durations)))
		fmt.Println("  - Median:", formatDuration(math.Median(summary.Durations)))
		fmt.Println("  - 90th percentile:", formatDuration(math.Percentile(summary.Durations, 0.9)))
		fmt.Println("  - 95th percentile:", formatDuration(math.Percentile(summary.Durations, 0.95)))
		fmt.Println("  - 99th percentile:", formatDuration(math.Percentile(summary.Durations, 0.99)))
	}
}

func printContributors(authors map[string]*analysis.RepositorySummary) {
	fmt.Println()
	fmt.Println("#### Issues")
	printIssueSummaryTable(authors)

	fmt.Println()
	fmt.Println("#### Pull Requests")
	printPullRequestSummaryTable(authors)
}

func printIssueSummaryTable(authors map[string]*analysis.RepositorySummary) {
	// sort by issue counts
	issueCounts := make(map[string]int)
	for author, summary := range authors {
		if summary.Issue.Total == 0 {
			continue
		}
		issueCounts[author] += summary.Issue.Total
	}
	counts := sort.SliceFromMap(issueCounts).Sort(func(s []sort.MapEntry[string, int], i, j int) bool {
		return s[i].Value > s[j].Value
	})

	// print header
	nameSize := 6 // len("Author")
	for _, entry := range counts {
		if len(entry.Key) > nameSize {
			nameSize = len(entry.Key)
		}
	}
	headerFormat := fmt.Sprintf("| %%-%ds | %%-6s | %%-6s | %%-6s | %%-8s | %%-8s | %%-8s | %%-8s | %%-8s |\n", nameSize)
	bodyFormat := fmt.Sprintf("| %%-%ds | %%-6d | %%-6d | %%-6d | %%-8s | %%-8s | %%-8s | %%-8s | %%-8s |\n", nameSize)
	fmt.Println()
	fmt.Printf(headerFormat, "Author", "Total", "Open", "Closed", "Min", "Max", "Mean", "Median", "P90")
	fmt.Printf("|%s|--------|--------|--------|----------|----------|----------|----------|----------|\n", strings.Repeat("-", nameSize+2))

	// print body
	for _, entry := range counts {
		author := entry.Key
		summary := authors[author].Issue
		if len(summary.Durations) == 0 {
			fmt.Printf(bodyFormat, author, summary.Total, summary.Open, summary.Closed, "", "", "", "", "")
			continue
		}

		sort.Sort(summary.Durations)
		fmt.Printf(bodyFormat, author, summary.Total, summary.Open, summary.Closed,
			formatDuration(math.Min(summary.Durations)),
			formatDuration(math.Max(summary.Durations)),
			formatDuration(math.Mean(summary.Durations)),
			formatDuration(math.Median(summary.Durations)),
			formatDuration(math.Percentile(summary.Durations, 0.9)),
		)
	}
}

func printPullRequestSummaryTable(authors map[string]*analysis.RepositorySummary) {
	// sort by pull request counts
	prCounts := make(map[string]int)
	for author, summary := range authors {
		if summary.PullRequest.Total == 0 {
			continue
		}
		prCounts[author] += summary.PullRequest.Total
	}
	counts := sort.SliceFromMap(prCounts).Sort(func(s []sort.MapEntry[string, int], i, j int) bool {
		return s[i].Value > s[j].Value
	})

	// print header
	nameSize := 6 // len("Author")
	for _, entry := range counts {
		if len(entry.Key) > nameSize {
			nameSize = len(entry.Key)
		}
	}
	headerFormat := fmt.Sprintf("| %%-%ds | %%-6s | %%-6s | %%-6s | %%-6s | %%-8s | %%-8s | %%-8s | %%-8s | %%-8s |\n", nameSize)
	bodyFormat := fmt.Sprintf("| %%-%ds | %%-6d | %%-6d | %%-6d | %%-6d | %%-8s | %%-8s | %%-8s | %%-8s | %%-8s |\n", nameSize)
	fmt.Println()
	fmt.Printf(headerFormat, "Author", "Total", "Open", "Closed", "Merged", "Min", "Max", "Mean", "Median", "P90")
	fmt.Printf("|%s|--------|--------|--------|--------|----------|----------|----------|----------|----------|\n", strings.Repeat("-", nameSize+2))

	// print body
	for _, entry := range counts {
		author := entry.Key
		summary := authors[author].PullRequest
		if len(summary.Durations) == 0 {
			fmt.Printf(bodyFormat, author, summary.Total, summary.Open, summary.Closed, summary.Merged, "", "", "", "", "")
			continue
		}

		sort.Sort(summary.Durations)
		fmt.Printf(bodyFormat, author, summary.Total, summary.Open, summary.Closed, summary.Merged,
			formatDuration(math.Min(summary.Durations)),
			formatDuration(math.Max(summary.Durations)),
			formatDuration(math.Mean(summary.Durations)),
			formatDuration(math.Median(summary.Durations)),
			formatDuration(math.Percentile(summary.Durations, 0.9)),
		)
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
