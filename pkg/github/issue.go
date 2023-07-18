package github

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shizhMSFT/gha/pkg/diff"
)

type Label struct {
	Name string `json:"name"`
}

func (l Label) String() string {
	return l.Name
}

type Account struct {
	Login string `json:"login"`
}

func (a Account) String() string {
	return a.Login
}

type Milestone struct {
	Title string `json:"title"`
}

func (m Milestone) String() string {
	return m.Title
}

type PullRequest struct {
	MergedAt *time.Time `json:"merged_at"`
}

// Issue is an abbreviated version of the GitHub issue type.
type Issue struct {
	HTMLURL     string       `json:"html_url"`
	Number      int          `json:"number"`
	Title       string       `json:"title"`
	User        Account      `json:"user"`
	Labels      []Label      `json:"labels"`
	Assignees   []Account    `json:"assignees"`
	State       string       `json:"state"`
	Milestone   Milestone    `json:"milestone"`
	CreatedAt   time.Time    `json:"created_at"`
	ClosedAt    *time.Time   `json:"closed_at"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`
}

func (i Issue) String() string {
	return fmt.Sprintf("%s #%d", i.Title, i.Number)
}

func (i Issue) IsPullRequest() bool {
	return i.PullRequest != nil
}

func (i Issue) Merged() bool {
	return i.IsPullRequest() && i.PullRequest.MergedAt != nil
}

func ParseIssues(jsonBytes []byte) (map[int]Issue, error) {
	var issues []Issue
	if err := json.Unmarshal(jsonBytes, &issues); err != nil {
		return nil, err
	}
	issueMap := make(map[int]Issue)
	for _, issue := range issues {
		issueMap[issue.Number] = issue
	}
	return issueMap, nil
}

// DiffIssues returns the difference between two sets of issues.
func DiffIssues(old, head map[int]Issue) map[int]diff.Diff[Issue] {
	diffs := make(map[int]diff.Diff[Issue])
	for number, prev := range old {
		curr, ok := head[number]
		if !ok {
			diffs[number] = diff.Diff[Issue]{
				Item: prev,
				Changes: []diff.Change{
					{Field: "State", Old: prev.State, New: "removed"},
				},
			}
			continue
		}

		var changes []diff.Change
		if change, changed := diff.DiffString("Title", prev.Title, curr.Title); changed {
			changes = append(changes, change)
		}
		if change, changed := diff.DiffSet("Labels", prev.Labels, curr.Labels); changed {
			changes = append(changes, change)
		}
		if change, changed := diff.DiffSet("Assignees", prev.Assignees, curr.Assignees); changed {
			changes = append(changes, change)
		}
		if change, changed := diff.DiffString("Milestone", prev.Milestone.String(), curr.Milestone.String()); changed {
			changes = append(changes, change)
		}
		if len(changes) > 0 {
			diffs[number] = diff.Diff[Issue]{Item: prev, Changes: changes}
		}
	}
	for number, curr := range head {
		if _, ok := old[number]; ok {
			continue
		}
		diffs[number] = diff.Diff[Issue]{
			Item: curr,
			Changes: []diff.Change{
				{Field: "State", Old: "new", New: curr.State},
			},
		}
	}
	return diffs
}
