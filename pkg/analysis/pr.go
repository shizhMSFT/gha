package analysis

import (
	"time"

	"github.com/shizhMSFT/gha/pkg/container/set"
	"github.com/shizhMSFT/gha/pkg/github"
)

type PullRequestReviewSummary struct {
	TimeFrame

	Reviewers map[string]set.Set[int]
}

func NewPullRequestReviewSummary() *PullRequestReviewSummary {
	return &PullRequestReviewSummary{
		Reviewers: make(map[string]set.Set[int]),
	}
}

func (s *PullRequestReviewSummary) ReviewCount() map[string]int {
	counts := make(map[string]int)
	for reviewer, numbers := range s.Reviewers {
		counts[reviewer] = numbers.Len()
	}
	return counts
}

func SummarizePullRequestReviews(reviews map[int][]github.PullRequestReview, start, end time.Time) *PullRequestReviewSummary {
	summary := NewPullRequestReviewSummary()
	summary.Start = start
	summary.End = end
	for number, reviews := range reviews {
		for _, review := range reviews {
			if (!start.IsZero() && review.SubmittedAt.Before(start)) || (!end.IsZero() && review.SubmittedAt.After(end)) {
				continue
			}
			reviewer := review.User.Login
			numbers := summary.Reviewers[reviewer]
			if numbers == nil {
				numbers = set.New[int]()
				summary.Reviewers[reviewer] = numbers
			}
			numbers.Add(number)
		}
	}
	return summary
}

type PullRequestReviewReport struct {
	TimeFrame

	Summaries map[string]*PullRequestReviewSummary
}

func (r *PullRequestReviewReport) Summarize(name string, reviews map[int][]github.PullRequestReview) *PullRequestReviewSummary {
	summary := SummarizePullRequestReviews(reviews, r.Start, r.End)
	r.Summaries[name] = summary
	return summary
}

func (r *PullRequestReviewReport) ReviewCount() map[string]int {
	counts := make(map[string]int)
	for _, summary := range r.Summaries {
		for reviewer, count := range summary.ReviewCount() {
			counts[reviewer] += count
		}
	}
	return counts
}
