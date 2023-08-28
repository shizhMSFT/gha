package analysis

import (
	"cmp"
	"slices"
	"strings"
	"time"

	"github.com/shizhMSFT/gha/pkg/container/set"
	"github.com/shizhMSFT/gha/pkg/github"
	"github.com/shizhMSFT/gha/pkg/sort"
)

type IssueCommentSummary struct {
	TimeFrame

	Responded  map[int]time.Duration // duration of first response
	NoResponse map[int]time.Time     // time of issue creation
}

func NewIssueCommentSummary() *IssueCommentSummary {
	return &IssueCommentSummary{
		Responded:  make(map[int]time.Duration),
		NoResponse: make(map[int]time.Time),
	}
}

type SummarizeIssueCommentsOptions struct {
	TimeFrame   TimeFrame
	Issues      map[int]github.Issue
	Comments    map[int][]github.IssueComment
	Maintainers set.Set[string]
}

func SummarizeIssueComments(opts SummarizeIssueCommentsOptions) *IssueCommentSummary {
	// normalize maintainers
	maintainers := set.New[string]()
	for maintainer := range opts.Maintainers {
		maintainers.Add(strings.ToLower(maintainer))
	}

	// summarize
	summary := NewIssueCommentSummary()
	summary.TimeFrame = opts.TimeFrame
	for number, comments := range opts.Comments {
		issue := opts.Issues[number]
		if issue.IsPullRequest() {
			continue
		}
		if !opts.TimeFrame.Contains(issue.CreatedAt) {
			continue
		}
		if maintainers.Contains(strings.ToLower(issue.User.Login)) {
			continue
		}
		slices.SortFunc(comments, func(a, b github.IssueComment) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
		responded := false
		for _, comment := range comments {
			if maintainers.Contains(strings.ToLower(comment.User.Login)) {
				duration := comment.CreatedAt.Sub(issue.CreatedAt)
				summary.Responded[number] = duration
				responded = true
				break
			}
		}
		if !responded {
			summary.NoResponse[number] = issue.CreatedAt
		}
	}
	return summary
}

func (s *IssueCommentSummary) OutOfSLA(sla time.Duration, now time.Time) sort.MapEntrySlice[int, time.Duration] {
	var durations sort.MapEntrySlice[int, time.Duration]
	for number, duration := range s.Responded {
		if duration > sla {
			durations = append(durations, sort.MapEntry[int, time.Duration]{
				Key:   number,
				Value: duration,
			})
		}
	}
	for number, createdAt := range s.NoResponse {
		duration := now.Sub(createdAt)
		if duration > sla {
			durations = append(durations, sort.MapEntry[int, time.Duration]{
				Key:   number,
				Value: duration,
			})
		}
	}
	durations.Sort(func(a, b sort.MapEntry[int, time.Duration]) int {
		return cmp.Compare(b.Value, a.Value)
	})
	return durations
}
