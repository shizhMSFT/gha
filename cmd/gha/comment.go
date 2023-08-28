package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/container/set"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/markdown"
	"github.com/shizhMSFT/gha/pkg/math"
	"github.com/urfave/cli/v3"
)

var maintainerRegexp = regexp.MustCompile(`(?:[^\w]+@|^@)([\w-]+)`)

var issueCommentCommand = &cli.Command{
	Name:      "issue-comment",
	Usage:     "analyze issue comments",
	ArgsUsage: "<issue_snapshot> <issue_comment_snapshot>",
	Aliases:   []string{"ic", "i"},
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
		&cli.IntFlag{
			Name:     "sla",
			Usage:    "report issues that have not received a comment from a maintainer more than `DAYS`",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "maintainers",
			Usage:    "specify a file containing a list of maintainer GitHub accounts",
			Required: true,
			OnlyOnce: true,
		},
	},
	Action: runIssueComment,
}

func runIssueComment(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		return errors.New("no issue or issue comment snapshot files specified")
	}

	// parse flags
	var opts analysis.SummarizeIssueCommentsOptions
	if ago := ctx.Int("ago"); ago > 0 {
		opts.TimeFrame.Start = time.Now().UTC().AddDate(0, 0, int(-ago))
	}
	if date := ctx.Value("start-date").(time.Time); !date.IsZero() {
		opts.TimeFrame.Start = date
	}
	if date := ctx.Value("end-date").(time.Time); !date.IsZero() {
		opts.TimeFrame.End = date
	}
	maintainers, err := readMaintainers(ctx.String("maintainers"))
	if err != nil {
		return err
	}
	opts.Maintainers = set.New(maintainers...)
	slaDays := ctx.Int("sla")
	sla := time.Duration(slaDays) * time.Hour * 24

	// read issue snapshot base
	opts.Issues, err = readIssues(ctx.Args().First())
	if err != nil {
		return err
	}

	// read issue comment snapshots
	opts.Comments, err = readIssueComments(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	// generate report
	fmt.Println("Issue Comment Summary")
	fmt.Println("=====================")
	printTimeFrame(opts.TimeFrame)
	fmt.Println()
	fmt.Println("## Maintainers")
	fmt.Println()
	for _, maintainer := range maintainers {
		fmt.Printf("- @%s\n", maintainer)
	}
	report := analysis.SummarizeIssueComments(opts)
	printIssueCommentSummary(report)
	if sla > 0 {
		fmt.Println()
		fmt.Println("### Out of SLA:", slaDays, "Days")
		fmt.Println()
		now := opts.TimeFrame.End
		if now.IsZero() {
			now = time.Now()
		}
		issues := report.OutOfSLA(sla, now)
		if len(issues) > 0 {
			table := markdown.NewTable("#Issue", "Duration", "Title")
			for _, issue := range issues {
				rawIssue := opts.Issues[issue.Key]
				table.AddRow(
					fmt.Sprintf("#%d", issue.Key),
					formatDuration(issue.Value),
					rawIssue.Title,
				)
			}
			table.Print(os.Stdout)
		} else {
			fmt.Println("No issues out of SLA")
		}
	}

	return nil
}

func readIssueComments(path string) (map[int][]github.IssueComment, error) {
	snapshotJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return github.ParseIssueComments(snapshotJSON)
}

func printIssueCommentSummary(summary *analysis.IssueCommentSummary) {
	fmt.Println()
	fmt.Println("## First Response Time")
	fmt.Println()
	fmt.Println("- Non-maintainer issues:", len(summary.Responded)+len(summary.NoResponse))
	fmt.Println("  - Responded:", len(summary.Responded))
	durations := make([]time.Duration, 0, len(summary.Responded))
	for _, duration := range summary.Responded {
		durations = append(durations, duration)
	}
	if len(durations) > 0 {
		slices.Sort(durations)
		fmt.Println("    - Min:", formatDuration(math.Min(durations)))
		fmt.Println("    - Max:", formatDuration(math.Max(durations)))
		fmt.Println("    - Mean:", formatDuration(math.Mean(durations)))
		fmt.Println("    - Median:", formatDuration(math.Median(durations)))
		fmt.Println("    - 90th percentile:", formatDuration(math.Percentile(durations, 0.9)))
		fmt.Println("    - 95th percentile:", formatDuration(math.Percentile(durations, 0.95)))
		fmt.Println("    - 99th percentile:", formatDuration(math.Percentile(durations, 0.99)))
	}
	fmt.Println("  - No Response:", len(summary.NoResponse))
}

func readMaintainers(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	matches := maintainerRegexp.FindAllSubmatch(content, -1)
	maintainers := make([]string, 0, len(matches))
	for _, match := range matches {
		maintainers = append(maintainers, string(match[1]))
	}
	return maintainers, nil
}
