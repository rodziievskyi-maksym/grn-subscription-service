package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/domain"
)

type SubscriptionRepositoryContract interface {
	GetUniqueRepositories(ctx context.Context) ([]string, error)
	GetOutdatedSubscriptions(ctx context.Context, repository, tag string) ([]domain.Subscription, error)
	GetSubscriptionsByEmail(ctx context.Context, email string) ([]domain.Subscription, error)

	CreateSubscription(ctx context.Context, subscription domain.Subscription) error
	UpdateLastTag(ctx context.Context, id uuid.UUID, tag string) error
	DeactivateSubscription(ctx context.Context, email, repo string) error
}
