package github

import (
	"encoding/json"
	"time"
)

type PullRequestReview struct {
	User        Account   `json:"user"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

func ParsePullRequestReviews(jsonBytes []byte) (map[int][]PullRequestReview, error) {
	var reviews map[int][]PullRequestReview
	if err := json.Unmarshal(jsonBytes, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}
