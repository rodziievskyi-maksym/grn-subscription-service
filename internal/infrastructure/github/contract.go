package github

import "context"

type Provider interface {
	GetLatestTag(ctx context.Context, owner, repo string) (string, error)
}
