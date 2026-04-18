package usecase

import (
	"context"

	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/domain"
)

type SubscriptionUseCaseContract interface {
	Subscribe(ctx context.Context, email, repository string) (*domain.Subscription, error)
	Unsubscribe(ctx context.Context, email, repository string) error
	GetSubscriptionsByEmail(ctx context.Context, email string) ([]domain.Subscription, error)
}
