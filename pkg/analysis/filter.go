package analysis

import "github.com/shizhMSFT/gha/pkg/github"

func Filter(issues map[int]github.Issue, filter func(github.Issue) bool) map[int]github.Issue {
	filtered := make(map[int]github.Issue)
	for number, issue := range issues {
		if filter(issue) {
			filtered[number] = issue
		}
	}
	return filtered
}
