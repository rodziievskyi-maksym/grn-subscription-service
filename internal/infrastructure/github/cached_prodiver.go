package github

import (
	"context"
	"fmt"

	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/infrastructure/cache"
)

type CachedGitHubProvider struct {
	realProvide Provider
	cache       *cache.TagCache
}

func NewCachedGitHubProvider(realProvide Provider, cache *cache.TagCache) *CachedGitHubProvider {
	return &CachedGitHubProvider{
		realProvide: realProvide,
		cache:       cache,
	}
}

// GetLatestTag is a decorator function the main aim is to get cached tags before reaching GitHub API
func (c *CachedGitHubProvider) GetLatestTag(ctx context.Context, owner, repo string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", owner, repo)

	// try to get tag from cache
	if tag, err := c.cache.GetTag(ctx, repoPath); err == nil {
		return tag, nil
	}

	// real call to GitHub API
	tag, err := c.realProvide.GetLatestTag(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	// set tag to cache
	if err := c.cache.SetTag(ctx, repoPath, tag); err != nil {
		return "", err
	}

	return tag, nil
}
