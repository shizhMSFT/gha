package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
func (c *Client) Snapshot(ctx context.Context, org, repo string) ([]byte, int, error) {
	var issues []json.RawMessage
	for page := 1; ; page++ {
		if c.PageEvent != nil {
			c.PageEvent(page)
		}
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all&direction=asc&per_page=100&page=%d", org, repo, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, 0, err
		}
		resp, err := c.do(req)
		if err != nil {
			return nil, 0, err
		}
		pagedIssues, err := func() ([]json.RawMessage, error) {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("%s: %s", url, resp.Status)
			}
			var items []json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return nil, err
			}
			return items, nil
		}()
		if err != nil {
			return nil, 0, err
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
