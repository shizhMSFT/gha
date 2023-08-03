package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shizhMSFT/gha/pkg/analysis"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/markdown"
	"github.com/shizhMSFT/gha/pkg/sort"
	"github.com/urfave/cli/v3"
)

var pullRequestReviewCommand = &cli.Command{
	Name:      "pr-review",
	Usage:     "analyze pull request reviews",
	ArgsUsage: "<review_snapshot> [...]",
	Aliases:   []string{"pr", "p"},
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
	counts := sort.SliceFromMap(reviewCounts).Sort(func(s []sort.MapEntry[string, int], i, j int) bool {
		return s[i].Value > s[j].Value
	})

	// print table
	table := markdown.NewTable("Reviewer", "Count", "")
	if len(counts) == 0 {
		return
	}
	barSize := 50
	max := counts[0].Value
	for _, entry := range counts {
		bar := strings.Repeat(" ", entry.Value*barSize/max)
		if bar == "" {
			bar = "` `"
		} else {
			bar = "`" + bar + "`"
		}
		table.AddRow(entry.Key, entry.Value, bar)
	}
	fmt.Println()
	table.Print(os.Stdout)
}
