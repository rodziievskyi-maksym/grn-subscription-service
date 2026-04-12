package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TagCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTagCache(client *redis.Client, ttl time.Duration) *TagCache {
	return &TagCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *TagCache) GetTag(ctx context.Context, repo string) (string, error) {
	return c.client.Get(ctx, repo).Result()
}

func (c *TagCache) SetTag(ctx context.Context, repo, tag string) error {
	return c.client.Set(ctx, repo, tag, c.ttl).Err()
}
