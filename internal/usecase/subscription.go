package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/domain"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/repository"
)

var (
	ErrInvalidRepositoryFormat    = errors.New("invalid repository format")
	ErrFailedToGetLatestTag       = errors.New("failed to get latest tag")
	ErrFailedToCreateSubscription = errors.New("failed to create subscription")
)

type SubscriptionUseCase struct {
	repository repository.SubscriptionRepositoryContract
	github     github.Provider
}

func NewSubscriptionUseCase(repo repository.SubscriptionRepositoryContract, gh github.Provider) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		repository: repo,
		github:     gh,
	}
}

func (s *SubscriptionUseCase) Subscribe(ctx context.Context, email, repository string) (*domain.Subscription, error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return nil, ErrInvalidRepositoryFormat
	}

	owner, repoName := parts[0], parts[1]

	tag, err := s.github.GetLatestTag(ctx, owner, repoName)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetLatestTag)
	}

	subscription := domain.NewSubscription(email, owner, repoName, tag)

	if err := s.repository.CreateSubscription(ctx, *subscription); err != nil {
		return nil, errors.Join(err, ErrFailedToCreateSubscription)
	}

	return subscription, nil
}

func (s *SubscriptionUseCase) Unsubscribe(ctx context.Context, email, repository string) error {
	if email == "" || repository == "" {
		return errors.New("email and repository are required")
	}

	return s.repository.DeactivateSubscription(ctx, email, repository)
}

func (s *SubscriptionUseCase) GetSubscriptionsByEmail(
	ctx context.Context,
	email string,
) ([]domain.Subscription, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	subs, err := s.repository.GetSubscriptionsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if subs == nil {
		return []domain.Subscription{}, nil
	}

	return subs, nil
}
