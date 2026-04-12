package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, addr, password string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.New("failed to ping redis")
	}

	return &Client{Client: client}, nil
}

func (r *Client) Close() {
	r.Client.Close()
}
