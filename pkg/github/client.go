package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// SnapshotOptions are options for taking a snapshot.
type SnapshotOptions struct {
	State        string
	UpdatedSince *time.Time
}

// Snapshot takes a snapshot of all issues and pull requests in a repository.
func (c *Client) Snapshot(ctx context.Context, org, repo string, opts SnapshotOptions) ([]byte, int, error) {
	switch opts.State = strings.ToLower(opts.State); opts.State {
	case "", "open", "closed", "all":
	default:
		return nil, 0, fmt.Errorf("invalid state: %s", opts.State)
	}
	var issues []json.RawMessage
	for page := 1; ; page++ {
		if c.PageEvent != nil {
			c.PageEvent(page)
		}
		u, err := url.Parse(fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", org, repo))
		if err != nil {
			return nil, 0, err
		}
		q := u.Query()
		q.Set("state", opts.State)
		q.Set("direction", "asc")
		q.Set("per_page", "100")
		q.Set("page", strconv.Itoa(page))
		if opts.UpdatedSince != nil {
			q.Set("since", opts.UpdatedSince.UTC().Format(time.RFC3339))
		}
		u.RawQuery = q.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, 0, err
		}
		resp, err := c.do(req)
		if err != nil {
			return nil, 0, err
		}
		pagedIssues, err := c.decodeResponse(resp)
		if err != nil {
			return nil, 0, fmt.Errorf("%v: %w", u, err)
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
