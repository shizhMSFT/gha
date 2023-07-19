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
	TimeFrame

	Authors map[string]*RepositorySummary
}

func NewSummary() *Summary {
	return &Summary{
		RepositorySummary: NewRepositorySummary(),
		Authors:           make(map[string]*RepositorySummary),
	}
}

func Summarize(issues map[int]github.Issue, start, end time.Time) *Summary {
	summary := NewSummary()
	summary.Start = start
	summary.End = end
	for _, issue := range issues {
		if (!start.IsZero() && issue.CreatedAt.Before(start)) || (!end.IsZero() && issue.CreatedAt.After(end)) {
			continue
		}
		author := issue.User.Login
		authorSummary := summary.Authors[author]
		if authorSummary == nil {
			authorSummary = NewRepositorySummary()
			summary.Authors[author] = authorSummary
		}
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
					duration := issue.ClosedAt.Sub(issue.CreatedAt)
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
				duration := issue.ClosedAt.Sub(issue.CreatedAt)
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

type Report struct {
	TimeFrame

	Summaries map[string]*Summary
}

func NewReport(start, end time.Time) *Report {
	return &Report{
		TimeFrame: TimeFrame{
			Start: start,
			End:   end,
		},
		Summaries: make(map[string]*Summary),
	}
}

func (r *Report) Summarize(name string, issues map[int]github.Issue) *Summary {
	summary := Summarize(issues, r.Start, r.End)
	r.Summaries[name] = summary
	return summary
}

func (r *Report) Abstract() *Summary {
	abstract := NewSummary()
	for _, summary := range r.Summaries {
		abstract.Combine(summary)
	}
	return abstract
}
