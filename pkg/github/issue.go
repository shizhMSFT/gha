package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
