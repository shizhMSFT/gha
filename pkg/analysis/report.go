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

func (s *IssueSummary) Union(other *IssueSummary) {
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

func (s *PullRequestSummary) Union(other *PullRequestSummary) {
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

func (s *RepositorySummary) Union(other *RepositorySummary) {
	s.Issue.Union(other.Issue)
	s.PullRequest.Union(other.PullRequest)
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

func Summarize(issues map[int]github.Issue, timeFrame TimeFrame) *Summary {
	summary := NewSummary()
	summary.TimeFrame = timeFrame
	for _, issue := range issues {
		if !timeFrame.Contains(issue.CreatedAt) {
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

func (s *Summary) Union(other *Summary) {
	s.RepositorySummary.Union(other.RepositorySummary)
	s.TimeFrame.Union(other.TimeFrame)

	for name, other := range other.Authors {
		summary := s.Authors[name]
		if summary == nil {
			summary = NewRepositorySummary()
			s.Authors[name] = summary
		}
		summary.Union(other)
	}
}

type Report struct {
	TimeFrame

	Summaries map[string]*Summary
}

func NewReport(timeFrame TimeFrame) *Report {
	return &Report{
		TimeFrame: timeFrame,
		Summaries: make(map[string]*Summary),
	}
}

func (r *Report) Summarize(name string, issues map[int]github.Issue) *Summary {
	summary := Summarize(issues, r.TimeFrame)
	r.Summaries[name] = summary
	return summary
}

func (r *Report) Abstract() *Summary {
	abstract := NewSummary()
	for _, summary := range r.Summaries {
		abstract.Union(summary)
	}
	return abstract
}
