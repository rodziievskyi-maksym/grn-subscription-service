package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/domain"
	"github.com/rodziievskyi-maksym/grn-subscription-service/pkg/databases/postgres"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	client *postgres.PostgreClient
}

func NewSubscriptionRepository(client *postgres.PostgreClient) SubscriptionRepositoryContract {
	return &SubscriptionRepository{client: client}
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, subscription domain.Subscription) error {
	query := `
        INSERT INTO subscriptions (id, email, repository, last_seen_tag, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (email, repository) 
        DO UPDATE SET 
            last_seen_tag = EXCLUDED.last_seen_tag, 
            updated_at = EXCLUDED.updated_at
        RETURNING id`

	err := r.client.Pool.QueryRow(ctx, query,
		subscription.ID,
		subscription.Email,
		subscription.Repository,
		subscription.LastSeenTag,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	).Scan(&subscription.ID)
	if err != nil {
		slog.Error("failed to create or update subscription", "error", err, "email", subscription.Email)
		return fmt.Errorf("failed to create subscription on data layer: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) GetUniqueRepositories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT repository FROM subscriptions WHERE is_active = true`

	rows, err := r.client.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query unique repositories: %w", err)
	}
	defer rows.Close()

	var repos []string

	for rows.Next() {
		var repo string
		if err := rows.Scan(&repo); err != nil {
			return nil, fmt.Errorf("failed to scan repository: %w", err)
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

func (r *SubscriptionRepository) GetOutdatedSubscriptions(
	ctx context.Context,
	repository, tag string,
) ([]domain.Subscription, error) {
	query := `
		SELECT id, email, repository, last_seen_tag 
		FROM subscriptions 
		WHERE repository = $1 AND last_seen_tag != $2`

	rows, err := r.client.Pool.Query(ctx, query, repository, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to query outdated subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.Email, &s.Repository, &s.LastSeenTag); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subs = append(subs, s)
	}

	return subs, nil
}

func (r *SubscriptionRepository) UpdateLastTag(ctx context.Context, id uuid.UUID, tag string) error {
	query := `UPDATE subscriptions SET last_seen_tag = $1 WHERE id = $2`

	result, err := r.client.Pool.Exec(ctx, query, tag, id)
	if err != nil {
		return fmt.Errorf("failed to update last seen tag: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no subscription found with id %s", id)
	}

	return nil
}

func (r *SubscriptionRepository) GetSubscriptionsByEmail(
	ctx context.Context,
	email string,
) ([]domain.Subscription, error) {
	query := `SELECT id, email, repository, last_seen_tag FROM subscriptions WHERE email = $1 AND is_active = true`

	rows, err := r.client.Pool.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("query subscriptions error: %w", err)
	}
	defer rows.Close()

	var subs []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.Email, &s.Repository, &s.LastSeenTag); err != nil {
			return nil, err
		}

		subs = append(subs, s)
	}

	return subs, nil
}

func (r *SubscriptionRepository) DeactivateSubscription(ctx context.Context, email, repo string) error {
	query := `UPDATE subscriptions SET is_active = false, updated_at = NOW() WHERE email = $1 AND repository = $2`

	res, err := r.client.Pool.Exec(ctx, query, email, repo)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}
