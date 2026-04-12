package usecase

import (
	"context"

	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/domain"
)

type SubscriptionUseCaseContract interface {
	Subscribe(ctx context.Context, email, repository string) (*domain.Subscription, error)
}
