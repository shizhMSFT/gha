package main

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/markdown"
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
		&cli.IntFlag{
			Name:     "issue-sla",
			Usage:    "report issues that were open for more than `DAYS`",
			OnlyOnce: true,
		},
		&cli.IntFlag{
			Name:     "pr-sla",
			Usage:    "report pull requests that were open for more than `DAYS`",
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
	opts := printSummaryOptions{
		includeContributors: ctx.Bool("contributors"),
		issueSLA:            int(ctx.Int("issue-sla")),
		pullRequestSLA:      int(ctx.Int("pr-sla")),
	}
	if ago := ctx.Int("ago"); ago > 0 {
		opts.timeFrame.Start = time.Now().UTC().AddDate(0, 0, int(-ago))
	}
	if date := ctx.Value("start-date").(time.Time); !date.IsZero() {
		opts.timeFrame.Start = date
	}
	if date := ctx.Value("end-date").(time.Time); !date.IsZero() {
		opts.timeFrame.End = date
	}

	// generate report
	fmt.Println("GitHub Analysis Report")
	fmt.Println("======================")
	printTimeFrame(opts.timeFrame)
	report := analysis.NewReport(opts.timeFrame)
	for _, path := range ctx.Args().Slice() {
		fmt.Println()
		fmt.Println("##", path)
		snapshot, err := readIssues(path)
		if err != nil {
			return err
		}
		printOpts := opts
		printOpts.snapshot = snapshot
		printSummary(report.Summarize(path, snapshot), printOpts)
	}
	if ctx.NArg() > 1 {
		fmt.Println()
		fmt.Println("## Overall")
		printSummary(report.Abstract(), opts)
	}
	return nil
}

func readIssues(path string) (map[int]github.Issue, error) {
	snapshotJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return github.ParseIssues(snapshotJSON)
}

func printTimeFrame(timeFrame analysis.TimeFrame) {
	if !timeFrame.Start.IsZero() {
		fmt.Printf("- Start Date: `%s`\n", timeFrame.Start.Format(time.DateTime))
	}
	if !timeFrame.End.IsZero() {
		fmt.Printf("- End Date: `%s`\n", timeFrame.End.Format(time.DateTime))
	}
}

type printSummaryOptions struct {
	timeFrame           analysis.TimeFrame
	includeContributors bool
	snapshot            map[int]github.Issue
	issueSLA            int
	pullRequestSLA      int
}

func printSummary(summary *analysis.Summary, opts printSummaryOptions) {
	printRepositorySummary(summary.RepositorySummary, opts)

	if opts.includeContributors {
		fmt.Println()
		fmt.Println("### Contributors")
		printContributors(summary.Authors)
	}
}

func printRepositorySummary(summary *analysis.RepositorySummary, opts printSummaryOptions) {
	fmt.Println()
	fmt.Println("### Issues")
	printIssueSummary(summary.Issue, opts)

	fmt.Println()
	fmt.Println("### Pull Requests")
	printPullRequestSummary(summary.PullRequest, opts)
}

func printIssueSummary(summary *analysis.IssueSummary, opts printSummaryOptions) {
	fmt.Println()
	fmt.Println("- Total:", summary.Total)
	fmt.Println("  - Open:", summary.Open)
	fmt.Println("  - Closed:", summary.Closed)
	if len(summary.Durations) > 0 {
		slices.Sort(summary.Durations)
		fmt.Println("- Time to close:")
		fmt.Println("  - Min:", formatDuration(math.Min(summary.Durations)))
		fmt.Println("  - Max:", formatDuration(math.Max(summary.Durations)))
		fmt.Println("  - Mean:", formatDuration(math.Mean(summary.Durations)))
		fmt.Println("  - Median:", formatDuration(math.Median(summary.Durations)))
		fmt.Println("  - 90th percentile:", formatDuration(math.Percentile(summary.Durations, 0.9)))
		fmt.Println("  - 95th percentile:", formatDuration(math.Percentile(summary.Durations, 0.95)))
		fmt.Println("  - 99th percentile:", formatDuration(math.Percentile(summary.Durations, 0.99)))
	}
	if opts.issueSLA > 0 && opts.snapshot != nil {
		fmt.Println()
		issues := analysis.Filter(opts.snapshot, func(issue github.Issue) bool {
			if issue.IsPullRequest() {
				return false
			}
			if !opts.timeFrame.Contains(issue.CreatedAt) {
				return false
			}
			if issue.State != "closed" {
				return false
			}
			duration := issue.Duration()
			return duration > time.Duration(opts.issueSLA)*time.Hour*24
		})
		if len(issues) == 0 {
			fmt.Println("No issues out of SLA:", opts.issueSLA, "days")
		} else {
			fmt.Println(len(issues), "issues out of SLA:", opts.issueSLA, "days")
			sortedIssues := sort.SliceFromMap(issues).Sort(func(a, b sort.MapEntry[int, github.Issue]) int {
				return cmp.Compare(b.Value.Duration(), a.Value.Duration())
			})
			table := markdown.NewTable("#Issue", "Duration", "Title")
			for _, issue := range sortedIssues {
				table.AddRow(
					fmt.Sprintf("#%d", issue.Value.Number),
					formatDuration(issue.Value.Duration()),
					issue.Value.Title,
				)
			}
			fmt.Println()
			table.Print(os.Stdout)
		}
	}
}

func printPullRequestSummary(summary *analysis.PullRequestSummary, opts printSummaryOptions) {
	fmt.Println()
	fmt.Println("- Total:", summary.Total)
	fmt.Println("  - Open:", summary.Open)
	fmt.Println("  - Closed:", summary.Closed)
	fmt.Println("  - Merged:", summary.Merged)
	if len(summary.Durations) > 0 {
		slices.Sort(summary.Durations)
		fmt.Println("- Time to merge:")
		fmt.Println("  - Min:", formatDuration(math.Min(summary.Durations)))
		fmt.Println("  - Max:", formatDuration(math.Max(summary.Durations)))
		fmt.Println("  - Mean:", formatDuration(math.Mean(summary.Durations)))
		fmt.Println("  - Median:", formatDuration(math.Median(summary.Durations)))
		fmt.Println("  - 90th percentile:", formatDuration(math.Percentile(summary.Durations, 0.9)))
		fmt.Println("  - 95th percentile:", formatDuration(math.Percentile(summary.Durations, 0.95)))
		fmt.Println("  - 99th percentile:", formatDuration(math.Percentile(summary.Durations, 0.99)))
	}
	if opts.pullRequestSLA > 0 && opts.snapshot != nil {
		fmt.Println()
		pullRequests := analysis.Filter(opts.snapshot, func(issue github.Issue) bool {
			if !issue.IsPullRequest() {
				return false
			}
			if !opts.timeFrame.Contains(issue.CreatedAt) {
				return false
			}
			if issue.State != "closed" {
				return false
			}
			if !issue.Merged() {
				return false
			}
			duration := issue.Duration()
			return duration > time.Duration(opts.pullRequestSLA)*time.Hour*24
		})
		if len(pullRequests) == 0 {
			fmt.Println("No pull requests out of SLA:", opts.pullRequestSLA, "days")
		} else {
			fmt.Println(len(pullRequests), "pull requests out of SLA:", opts.pullRequestSLA, "days")
			sortedPullRequests := sort.SliceFromMap(pullRequests).Sort(func(a, b sort.MapEntry[int, github.Issue]) int {
				return cmp.Compare(b.Value.Duration(), a.Value.Duration())
			})
			table := markdown.NewTable("#PR", "Duration", "Title")
			for _, pullRequest := range sortedPullRequests {
				table.AddRow(pullRequest.Value.Number, formatDuration(pullRequest.Value.Duration()), pullRequest.Value.Title)
			}
			fmt.Println()
			table.Print(os.Stdout)
		}
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
	counts := sort.SliceFromMap(issueCounts).Sort(func(a, b sort.MapEntry[string, int]) int {
		return b.Value - a.Value
	})

	// print table
	table := markdown.NewTable("Author", "Total", "Open", "Closed", "Min", "Max", "Mean", "Median", "P90")
	for _, entry := range counts {
		author := entry.Key
		summary := authors[author].Issue
		if len(summary.Durations) == 0 {
			table.AddRow(author, summary.Total, summary.Open, summary.Closed, "-", "-", "-", "-", "-")
			continue
		}

		slices.Sort(summary.Durations)
		table.AddRow(author, summary.Total, summary.Open, summary.Closed,
			formatDuration(math.Min(summary.Durations)),
			formatDuration(math.Max(summary.Durations)),
			formatDuration(math.Mean(summary.Durations)),
			formatDuration(math.Median(summary.Durations)),
			formatDuration(math.Percentile(summary.Durations, 0.9)),
		)
	}
	fmt.Println()
	table.Print(os.Stdout)
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
	counts := sort.SliceFromMap(prCounts).Sort(func(a, b sort.MapEntry[string, int]) int {
		return b.Value - a.Value
	})

	// print table
	table := markdown.NewTable("Author", "Total", "Open", "Closed", "Merged", "Min", "Max", "Mean", "Median", "P90")
	for _, entry := range counts {
		author := entry.Key
		summary := authors[author].PullRequest
		if len(summary.Durations) == 0 {
			table.AddRow(author, summary.Total, summary.Open, summary.Closed, summary.Merged, "-", "-", "-", "-", "-")
			continue
		}

		slices.Sort(summary.Durations)
		table.AddRow(author, summary.Total, summary.Open, summary.Closed, summary.Merged,
			formatDuration(math.Min(summary.Durations)),
			formatDuration(math.Max(summary.Durations)),
			formatDuration(math.Mean(summary.Durations)),
			formatDuration(math.Median(summary.Durations)),
			formatDuration(math.Percentile(summary.Durations, 0.9)),
		)
	}
	fmt.Println()
	table.Print(os.Stdout)
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
