package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/domain"
)

type SubscriptionRepositoryContract interface {
	CreateSubscription(ctx context.Context, subscription domain.Subscription) error
	GetUniqueRepositories(ctx context.Context) ([]string, error)
	GetOutdatedSubscriptions(ctx context.Context, repository, tag string) ([]domain.Subscription, error)
	UpdateLastTag(ctx context.Context, id uuid.UUID, tag string) error
}
