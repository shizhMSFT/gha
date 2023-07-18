package analysis

import (
	"time"

	"github.com/shizhMSFT/gha/pkg/github"
)

type Report struct {
	Start time.Time
	End   time.Time

	Summaries map[string]*Summary
}

func NewReport(start, end time.Time) *Report {
	return &Report{
		Start:     start,
		End:       end,
		Summaries: make(map[string]*Summary),
	}
}

func (r *Report) Summarize(name string, issues map[int]github.Issue) *Summary {
	summary := Summarize(issues, r.Start, r.End)
	r.Summaries[name] = summary
	return summary
}

func (r *Report) Abstract() *Summary {
	abstract := new(Summary)
	for _, summary := range r.Summaries {
		abstract.Combine(summary)
	}
	return abstract
}
