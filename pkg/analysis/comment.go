package analysis

import (
	"slices"
	"time"

	"github.com/shizhMSFT/gha/pkg/container/set"
	"github.com/shizhMSFT/gha/pkg/github"
)

type IssueCommentSummary struct {
	TimeFrame

	Responded  map[int]time.Duration // duration of first response
	NoResponse set.Set[int]
}

func NewIssueCommentSummary() *IssueCommentSummary {
	return &IssueCommentSummary{
		Responded:  make(map[int]time.Duration),
		NoResponse: set.New[int](),
	}
}

type SummarizeIssueCommentsOptions struct {
	TimeFrame   TimeFrame
	Issues      map[int]github.Issue
	Comments    map[int][]github.IssueComment
	Maintainers set.Set[string]
}

func SummarizeIssueComments(opts SummarizeIssueCommentsOptions) *IssueCommentSummary {
	summary := NewIssueCommentSummary()
	summary.TimeFrame = opts.TimeFrame
	for number, comments := range opts.Comments {
		issue := opts.Issues[number]
		if !opts.TimeFrame.Contains(issue.CreatedAt) {
			continue
		}
		if opts.Maintainers.Contains(issue.User.Login) {
			continue
		}
		slices.SortFunc(comments, func(a, b github.IssueComment) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
		responded := false
		for _, comment := range comments {
			if opts.Maintainers.Contains(comment.User.Login) {
				duration := comment.CreatedAt.Sub(issue.CreatedAt)
				summary.Responded[number] = duration
				responded = true
				break
			}
		}
		if !responded {
			summary.NoResponse.Add(number)
		}
	}
	return nil
}
