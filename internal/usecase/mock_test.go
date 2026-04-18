package usecase_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetOutdatedSubscriptions(ctx context.Context, repository, tag string) ([]domain.Subscription, error) {
	args := m.Called(ctx, repository, tag)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.Subscription), args.Error(1)
}

func (m *MockRepository) GetSubscriptionsByEmail(ctx context.Context, email string) ([]domain.Subscription, error) {
	args := m.Called(ctx, email)
	return args.Get(0).([]domain.Subscription), args.Error(1)
}

func (m *MockRepository) GetUniqueRepositories(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) DeactivateSubscription(ctx context.Context, email, repo string) error {
	args := m.Called(ctx, email, repo)
	return args.Error(0)
}

func (m *MockRepository) CreateSubscription(ctx context.Context, subscription domain.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockRepository) UpdateLastTag(ctx context.Context, id uuid.UUID, tag string) error {
	args := m.Called(ctx, id, tag)
	return args.Error(0)
}

type MockGitHub struct {
	mock.Mock
}

func (m *MockGitHub) GetLatestTag(ctx context.Context, owner, repo string) (string, error) {
	args := m.Called(ctx, owner, repo)
	return args.String(0), args.Error(1)
}
