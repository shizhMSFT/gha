package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/sort"
	"github.com/urfave/cli/v3"
)

var pullRequestReviewCommand = &cli.Command{
	Name:      "pull-request-review",
	Usage:     "analyze pull request reviews",
	ArgsUsage: "<review_snapshot> [...]",
	Aliases:   []string{"pr-review", "prr", "pr", "p"},
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
	Action: runPullRequestReview,
}

func runPullRequestReview(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return errors.New("no review snapshot files specified")
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

	// generate report
	fmt.Println("Pull Request Review Count")
	fmt.Println("==========================")
	printTimeFrame(timeFrame)
	report := analysis.NewPullRequestReviewReport(timeFrame)
	for _, path := range ctx.Args().Slice() {
		fmt.Println()
		fmt.Println("##", path)
		snapshotJSON, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot, err := github.ParsePullRequestReviews(snapshotJSON)
		if err != nil {
			return err
		}
		printPullRequestReviewCount(report.Summarize(path, snapshot).ReviewCount())
	}
	if ctx.NArg() > 1 {
		fmt.Println()
		fmt.Println("## Overall")
		printPullRequestReviewCount(report.ReviewCount())
	}
	return nil
}

func printPullRequestReviewCount(reviewCounts map[string]int) {
	// sort by review counts
	counts := sort.SliceFromMap(reviewCounts, false)

	// print header
	nameSize := 8 // len("Reviewer")
	for _, entry := range counts {
		if len(entry.Key) > nameSize {
			nameSize = len(entry.Key)
		}
	}
	barSize := 50
	headerFormat := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds |\n", nameSize, 5, barSize+2)
	bodyFormat := fmt.Sprintf("| %%-%ds | %%%dd | %%-%ds |\n", nameSize, 5, barSize+2)
	fmt.Println()
	fmt.Printf(headerFormat, "Reviewer", "Count", "")
	fmt.Printf("|%s|%s|%s|\n", strings.Repeat("-", nameSize+2), strings.Repeat("-", 7), strings.Repeat("-", barSize+4))

	// print body
	if len(counts) == 0 {
		return
	}
	max := counts[0].Value
	for _, entry := range counts {
		bar := strings.Repeat(" ", entry.Value*barSize/max)
		if bar == "" {
			bar = "` `"
		} else {
			bar = "`" + bar + "`"
		}
		fmt.Printf(bodyFormat, entry.Key, entry.Value, bar)
	}
}
