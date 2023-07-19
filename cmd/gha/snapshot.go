package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/urfave/cli/v3"
)

var snapshotCommand = &cli.Command{
	Name:      "snapshot",
	ArgsUsage: "<org>/<repo>",
	Usage:     "take a snapshot of a repository",
	Aliases:   []string{"s"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "pr-reviews",
			Usage:    "include pull request reviews in the snapshot",
			OnlyOnce: true,
		},
		&cli.IntFlag{
			Name:     "pr-reviews-ago",
			Usage:    "include pull request reviews since `DAYS` ago",
			OnlyOnce: true,
		},
		&cli.TimestampFlag{
			Name:     "pr-reviews-since",
			Usage:    "include pull request reviews since `DATE`",
			Config:   cli.TimestampConfig{Layout: time.DateOnly},
			OnlyOnce: true,
		},
	},
	Action: runSnapshot,
}

func runSnapshot(ctx *cli.Context) error {
	ref := ctx.Args().First()
	org, repo, ok := strings.Cut(ref, "/")
	if !ok {
		return fmt.Errorf("invalid ref: %s", ref)
	}

	client := github.NewClient()
	client.Token = os.Getenv("GITHUB_TOKEN")
	client.PageEvent = func(page int) {
		fmt.Printf(".")
	}
	snapshot, n, err := client.Snapshot(ctx.Context, org, repo)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("Fetched", n, "issues and pull requests")

	path := fmt.Sprintf("%s_%s_%s_snapshot.json", org, repo, time.Now().Format("20060102_150405"))
	if err := os.WriteFile(path, snapshot, 0644); err != nil {
		return err
	}
	fmt.Println("Saved snapshot to", path)

	if ctx.Bool("pr-reviews") {
		if err := snapshotPullRequestReviews(ctx, org, repo, client, snapshot); err != nil {
			return err
		}
	}

	return nil
}

func snapshotPullRequestReviews(ctx *cli.Context, org string, repo string, client *github.Client, snapshot []byte) error {
	// parse flags
	var start time.Time
	if ago := ctx.Int("pr-reviews-ago"); ago > 0 {
		start = time.Now().UTC().AddDate(0, 0, int(-ago))
	}
	if date := ctx.Value("pr-reviews-since").(time.Time); !date.IsZero() {
		start = date
	}

	// fetch pull request reviews
	issues, err := github.ParseIssues(snapshot)
	if err != nil {
		return err
	}
	reviews := make(map[int]json.RawMessage)
	for _, issue := range issues {
		if issue.IsPullRequest() {
			if !start.IsZero() && issue.CreatedAt.Before(start) {
				continue
			}
			reviews[issue.Number] = nil
		}
	}
	total := len(reviews)
	fmt.Printf("Fetching reviews of %d pull requests", total)
	if !start.IsZero() {
		fmt.Printf(" since %s", start.Format(time.DateOnly))
	}
	fmt.Println("...")

	count := 0
	for number := range reviews {
		review, err := client.PullRequestReviews(ctx.Context, org, repo, number)
		if err != nil {
			return err
		}
		reviews[number] = review
		count++
		fmt.Printf(".")
		if count%50 == 0 {
			fmt.Printf(" %6g%%\n", float64(10000*count/total)/100.0)
		}
	}
	fmt.Println(strings.Repeat(" ", 50-count%50), "100.00%")

	// save reviews
	path := fmt.Sprintf("%s_%s_%s_reviews.json", org, repo, time.Now().Format("20060102_150405"))
	reviewsJSON, err := json.Marshal(reviews)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, reviewsJSON, 0644); err != nil {
		return err
	}
	fmt.Println("Saved pull request reviews to", path)

	return nil
}
