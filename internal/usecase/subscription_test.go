package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/domain"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/infrastructure/github"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/infrastructure/repository"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionUseCase_Subscribe(t *testing.T) {
	type testCase struct {
		name        string
		email       string
		repoPath    string
		mockGHRes   string
		mockGHErr   error
		mockRepoErr error
		expectedErr error
		expectTag   string
	}

	tests := []testCase{
		{
			name:        "Success - New subscription",
			email:       "test@test.com",
			repoPath:    "owner/repo",
			mockGHRes:   "v1.0.0",
			mockGHErr:   nil,
			mockRepoErr: nil,
			expectedErr: nil,
			expectTag:   "v1.0.0",
		},
		{
			name:        "Failure - Repository Not Found",
			email:       "test@test.com",
			repoPath:    "invalid/repo",
			mockGHRes:   "",
			mockGHErr:   github.ErrRepositoryNotFound,
			mockRepoErr: nil,
			expectedErr: github.ErrRepositoryNotFound,
			expectTag:   "",
		},
		{
			name:        "Failure - Database Error",
			email:       "test@test.com",
			repoPath:    "owner/repo",
			mockGHRes:   "v1.0.0",
			mockGHErr:   nil,
			mockRepoErr: errors.New("db connection lost"),
			expectedErr: errors.New("db connection lost"),
			expectTag:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockGH := new(MockGitHub)
			uc := usecase.NewSubscriptionUseCase(mockRepo, mockGH)
			ctx := context.Background()

			if tc.repoPath != "" {
				owner, repo, _ := strings.Cut(tc.repoPath, "/")
				mockGH.On("GetLatestTag", ctx, owner, repo).Return(tc.mockGHRes, tc.mockGHErr)
			}

			if tc.mockGHErr == nil && tc.repoPath != "" {
				mockRepo.On("CreateSubscription", ctx, mock.AnythingOfType("domain.Subscription")).
					Return(tc.mockRepoErr)
			}

			res, err := uc.Subscribe(ctx, tc.email, tc.repoPath)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectTag, res.LastSeenTag)
			}
		})
	}
}

func TestSubscriptionUseCase_Unsubscribe(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		repo        string
		mockErr     error
		expectedErr error
	}{
		{
			name:        "Success - Deactivate active sub",
			email:       "user@test.com",
			repo:        "owner/repo",
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "Failure - Subscription not found",
			email:       "user@test.com",
			repo:        "owner/repo",
			mockErr:     repository.ErrSubscriptionNotFound,
			expectedErr: repository.ErrSubscriptionNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			uc := usecase.NewSubscriptionUseCase(mockRepo, nil)

			mockRepo.On("DeactivateSubscription", mock.Anything, tc.email, tc.repo).
				Return(tc.mockErr)

			err := uc.Unsubscribe(context.Background(), tc.email, tc.repo)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSubscriptionUseCase_GetSubscriptions(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		mockReturn   []domain.Subscription
		mockErr      error
		expectedLen  int
		expectNilRes bool
	}{
		{
			name:  "Success - Found 2 subscriptions",
			email: "user@test.com",
			mockReturn: []domain.Subscription{
				{Repository: "owner/repo1"},
				{Repository: "owner/repo2"},
			},
			mockErr:     nil,
			expectedLen: 2,
		},
		{
			name:         "Success - No subscriptions found (Empty Slice)",
			email:        "new@user.com",
			mockReturn:   []domain.Subscription{},
			mockErr:      nil,
			expectedLen:  0,
			expectNilRes: false,
		},
		{
			name:        "Failure - Database Error",
			email:       "error@test.com",
			mockReturn:  nil,
			mockErr:     errors.New("db error"),
			expectedLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			uc := usecase.NewSubscriptionUseCase(mockRepo, nil)

			mockRepo.On("GetSubscriptionsByEmail", mock.Anything, tc.email).
				Return(tc.mockReturn, tc.mockErr)

			res, err := uc.GetSubscriptionsByEmail(context.Background(), tc.email)

			if tc.mockErr != nil {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, res, tc.expectedLen)
				assert.NotNil(t, res)
			}
		})
	}
}
