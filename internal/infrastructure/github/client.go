package github

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/go-github/v84/github"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
)

type RateLimitError struct {
	ResetTime time.Time
	Message   string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("%s: resets at %s", e.Message, e.ResetTime.Format(time.RFC3339))
}

func IsRateLimitError(err error) bool {
	var rlErr *RateLimitError
	return errors.As(err, &rlErr)
}

type client struct {
	api *github.Client
}

// NewClient creates a new GitHub client.
func NewClient(token string) Provider {
	return &client{api: github.NewClient(nil).WithAuthToken(token)}
}

func (c *client) GetLatestTag(ctx context.Context, owner, repo string) (string, error) {
	release, resp, err := c.api.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusNotFound:
				return "", ErrRepositoryNotFound
			case http.StatusTooManyRequests, http.StatusForbidden:
				if resp.Rate.Remaining == 0 {
					return "", &RateLimitError{
						ResetTime: resp.Rate.Reset.Time,
						Message:   resp.Status,
					}
				}
			}
		}

		slog.Error("github api error", "owner", owner, "repo", repo, "err", err)

		return "", err
	}

	if release == nil || release.TagName == nil {
		return "", errors.New("no releases found for repository")
	}

	return *release.TagName, nil
}
