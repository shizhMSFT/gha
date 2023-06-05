package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shizhMSFT/ghutil/pkg/diff"
)

// OpenIssues returns all issues that are still open.
func OpenIssues(ctx context.Context, org, repo string) ([]json.RawMessage, error) {
	var openIssues []json.RawMessage
	for page := 1; ; page++ {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?per_page=100&page=%d", org, repo, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-Github-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: %s", url, resp.Status)
		}

		var issues []json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
			return nil, err
		}
		openIssues = append(openIssues, issues...)
		if len(issues) < 100 {
			break
		}
	}
	return openIssues, nil
}

type Label struct {
	Name string `json:"name"`
}

func (l Label) String() string {
	return l.Name
}

type Assignee struct {
	Login string `json:"login"`
}

func (a Assignee) String() string {
	return a.Login
}

type Milestone struct {
	Title string `json:"title"`
}

func (m Milestone) String() string {
	return m.Title
}

// Issue is an abbreviated version of the GitHub issue type.
type Issue struct {
	HTMLURL   string     `json:"html_url"`
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Labels    []Label    `json:"labels"`
	Assignees []Assignee `json:"assignees"`
	State     string     `json:"state"`
	Milestone Milestone  `json:"milestone"`
}

// DiffIssues returns the difference between two sets of issues.
func DiffIssues(old, new []byte) ([]diff.Diff[Issue], error) {
	oldSet, err := issuesToSet(old)
	if err != nil {
		return nil, err
	}
	newSet, err := issuesToSet(new)
	if err != nil {
		return nil, err
	}

	var diffs []diff.Diff[Issue]
	for number, oldIssue := range oldSet {
		newIssue, ok := newSet[number]
		if !ok {
			diffs = append(diffs, diff.Diff[Issue]{
				Item: oldIssue,
				Changes: []diff.Change{
					{Field: "State", Old: oldIssue.State, New: "closed"},
				},
			})
			continue
		}

		var changes []diff.Change
		if change, changed := diff.DiffString("Title", oldIssue.Title, newIssue.Title); changed {
			changes = append(changes, change)
		}
		if change, changed := diff.DiffSet("Labels", oldIssue.Labels, newIssue.Labels); changed {
			changes = append(changes, change)
		}
		if change, changed := diff.DiffSet("Assignees", oldIssue.Assignees, newIssue.Assignees); changed {
			changes = append(changes, change)
		}
		if change, changed := diff.DiffString("Milestone", oldIssue.Milestone.String(), newIssue.Milestone.String()); changed {
			changes = append(changes, change)
		}
		if len(changes) > 0 {
			diffs = append(diffs, diff.Diff[Issue]{Item: oldIssue, Changes: changes})
		}
		delete(newSet, number)
	}
	for _, newIssue := range newSet {
		diffs = append(diffs, diff.Diff[Issue]{
			Item: newIssue,
			Changes: []diff.Change{
				{Field: "State", Old: "closed or new", New: newIssue.State},
			},
		})
	}
	return diffs, nil
}

func issuesToSet(jsonBytes []byte) (map[int]Issue, error) {
	var issues []Issue
	if err := json.Unmarshal(jsonBytes, &issues); err != nil {
		return nil, err
	}
	issueSet := make(map[int]Issue)
	for _, issue := range issues {
		issueSet[issue.Number] = issue
	}
	return issueSet, nil
}
