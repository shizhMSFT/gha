package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Client     *http.Client
	APIVersion string
	Token      string
	PageEvent  func(page int) // PageEvent is called when a new page is being fetched.
}

func NewClient() *Client {
	return &Client{
		Client:     http.DefaultClient,
		APIVersion: "2022-11-28",
	}
}

// Snapshot takes a snapshot of all issues and pull requests in a repository.
func (c *Client) Snapshot(ctx context.Context, org, repo string, updatedSince time.Time) ([]byte, int, error) {
	var issues []json.RawMessage
	for page := 1; ; page++ {
		if c.PageEvent != nil {
			c.PageEvent(page)
		}
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all&direction=asc&per_page=100&page=%d", org, repo, page)
		if !updatedSince.IsZero() {
			url += "&since=" + updatedSince.UTC().Format(time.RFC3339)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, 0, err
		}
		resp, err := c.do(req)
		if err != nil {
			return nil, 0, err
		}
		pagedIssues, err := c.decodeResponse(resp)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", url, err)
		}
		issues = append(issues, pagedIssues...)
		if len(pagedIssues) < 100 {
			break
		}
	}
	snapshot, err := json.Marshal(issues)
	if err != nil {
		return nil, 0, err
	}
	return snapshot, len(issues), nil
}

// PullRequestReviews takes a snapshot of all reviews for a pull request.
func (c *Client) PullRequestReviews(ctx context.Context, org, repo string, number int) ([]byte, error) {
	var reviews []json.RawMessage
	for page := 1; ; page++ {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/reviews?per_page=100&page=%d", org, repo, number, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.do(req)
		if err != nil {
			return nil, err
		}
		pagedReviews, err := c.decodeResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", url, err)
		}
		reviews = append(reviews, pagedReviews...)
		if len(pagedReviews) < 100 {
			break
		}
	}
	return json.Marshal(reviews)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/vnd.github+json")
	apiVersion := c.APIVersion
	if apiVersion == "" {
		apiVersion = "2022-11-28"
	}
	req.Header.Set("X-Github-Api-Version", apiVersion)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return c.Client.Do(req)
}

func (c *Client) decodeResponse(resp *http.Response) ([]json.RawMessage, error) {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden && c.Token == "" {
			return nil, fmt.Errorf("%s: provide GITHUB_TOKEN may help", resp.Status)
		}
		return nil, errors.New(resp.Status)
	}
	var items []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
