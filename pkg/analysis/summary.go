package analysis

import (
	"time"

	"github.com/shizhMSFT/gha/pkg/github"
)

type IssueSummary struct {
	Total     int
	Open      int
	Closed    int
	Durations []time.Duration // time to close
}

func (s *IssueSummary) Combine(other *IssueSummary) {
	s.Total += other.Total
	s.Open += other.Open
	s.Closed += other.Closed
	s.Durations = append(s.Durations, other.Durations...)
}

type PullRequestSummary struct {
	Total     int
	Open      int
	Closed    int
	Merged    int
	Durations []time.Duration // time to merge
}

func (s *PullRequestSummary) Combine(other *PullRequestSummary) {
	s.Total += other.Total
	s.Open += other.Open
	s.Closed += other.Closed
	s.Merged += other.Merged
	s.Durations = append(s.Durations, other.Durations...)
}

type RepositorySummary struct {
	Issue       *IssueSummary
	PullRequest *PullRequestSummary
}

func NewRepositorySummary() *RepositorySummary {
	return &RepositorySummary{
		Issue:       new(IssueSummary),
		PullRequest: new(PullRequestSummary),
	}
}

func (s *RepositorySummary) Combine(other *RepositorySummary) {
	s.Issue.Combine(other.Issue)
	s.PullRequest.Combine(other.PullRequest)
}

type Summary struct {
	*RepositorySummary

	Start time.Time
	End   time.Time

	Authors map[string]*RepositorySummary
}

func Summarize(issues map[int]github.Issue, start, end time.Time) *Summary {
	summary := &Summary{
		RepositorySummary: NewRepositorySummary(),
		Start:             start,
		End:               end,
		Authors:           make(map[string]*RepositorySummary),
	}
	for _, issue := range issues {
		if (!issue.CreatedAt.IsZero() && issue.CreatedAt.Before(start)) || (!issue.ClosedAt.IsZero() && issue.ClosedAt.After(end)) {
			continue
		}
		author := issue.User.Login
		authorSummary := summary.Authors[author]
		if authorSummary == nil {
			authorSummary = NewRepositorySummary()
			summary.Authors[author] = authorSummary
		}
		duration := issue.ClosedAt.Sub(issue.CreatedAt)
		if issue.IsPullRequest() {
			summary.PullRequest.Total++
			authorSummary.PullRequest.Total++
			switch issue.State {
			case "open":
				summary.PullRequest.Open++
				authorSummary.PullRequest.Open++
			case "closed":
				if issue.Merged() {
					summary.PullRequest.Merged++
					authorSummary.PullRequest.Merged++
					summary.PullRequest.Durations = append(summary.PullRequest.Durations, duration)
					authorSummary.PullRequest.Durations = append(authorSummary.PullRequest.Durations, duration)
				} else {
					summary.PullRequest.Closed++
					authorSummary.PullRequest.Closed++
				}
			}
		} else {
			summary.Issue.Total++
			authorSummary.Issue.Total++
			switch issue.State {
			case "open":
				summary.Issue.Open++
				authorSummary.Issue.Open++
			case "closed":
				summary.Issue.Closed++
				authorSummary.Issue.Closed++
				summary.Issue.Durations = append(summary.Issue.Durations, duration)
				authorSummary.Issue.Durations = append(authorSummary.Issue.Durations, duration)
			}
		}
	}
	return summary
}

func (s *Summary) Combine(other *Summary) {
	s.RepositorySummary.Combine(other.RepositorySummary)

	if s.Start.After(other.Start) {
		s.Start = other.Start
	}
	if s.End.Before(other.End) {
		s.End = other.End
	}

	for name, other := range other.Authors {
		summary := s.Authors[name]
		if summary == nil {
			summary = NewRepositorySummary()
			s.Authors[name] = summary
		}
		summary.Combine(other)
	}
}
