package worker

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/config"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/email"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/repository"
)

type Scanner struct {
	repository    repository.SubscriptionRepositoryContract
	github        github.Provider
	emailProvider email.EmailProvider
	scheduler     gocron.Scheduler

	rateLimitUntil time.Time
	mu             sync.RWMutex
}

func NewScanner(repo repository.SubscriptionRepositoryContract, gh github.Provider) (*Scanner, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, errors.Join(errors.New("failed to create scheduler"), err)
	}

	return &Scanner{
		repository:    repo,
		github:        gh,
		scheduler:     scheduler,
		emailProvider: email.NewSmtpProvider(),
	}, nil
}

func (s *Scanner) Run(ctx context.Context) error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(config.Cfg().ScannerInterval),
		gocron.NewTask(s.checkUpdates, ctx),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.JobOption(gocron.WithStartImmediately()),
	)
	if err != nil {
		return errors.Join(errors.New("failed to create job"), err)
	}

	s.scheduler.Start()

	<-ctx.Done()

	return s.scheduler.Shutdown()
}

func (s *Scanner) checkUpdates(ctx context.Context) {
	s.mu.Lock()
	until := s.rateLimitUntil
	defer s.mu.Unlock()

	if time.Now().Before(until) {
		slog.Warn("Scanner: skipping check due to active rate limit",
			"until", until.Format(time.RFC3339))
		return
	}

	slog.Info("Scanner: checking updates...")

	//get unique repositories
	repositories, err := s.repository.GetUniqueRepositories(ctx)
	if err != nil {
		slog.Error("Scanner: failed to get unique repositories", "error", err)
		return
	}

	for _, repoPath := range repositories {
		parts := strings.Split(repoPath, "/")
		if len(parts) != 2 {
			slog.Warn("Scanner: Invalid repository format in DB", "path", repoPath)
			continue
		}

		newTag, err := s.github.GetLatestTag(ctx, parts[0], parts[1])
		if err != nil {
			if github.IsRateLimitError(err) {
				var rlErr *github.RateLimitError
				errors.As(err, &rlErr)

				s.mu.Lock()
				s.rateLimitUntil = rlErr.ResetTime
				s.mu.Unlock()

				slog.Warn("Scanner: pausing due to rate limit", "until", rlErr.ResetTime)
				return
			}
			slog.Error("Scanner: GitHub API error", "repo", repoPath, "error", err)
			continue
		}

		outdatedSubs, err := s.repository.GetOutdatedSubscriptions(ctx, repoPath, newTag)
		if err != nil {
			slog.Error("Scanner: Failed to fetch outdated subscriptions", "repo", repoPath, "error", err)
			continue
		}

		if len(outdatedSubs) == 0 {
			slog.Debug("Scanner: No new updates for repository", "repo", repoPath)
			continue
		}

		slog.Info("Scanner: Found updates", "repo", repoPath, "count", len(outdatedSubs), "new_tag", newTag)

		for _, sub := range outdatedSubs {
			err := s.emailProvider.SendReleaseNotification(ctx, sub.Email, sub.Repository, newTag)
			if err != nil {
				slog.Error("Failed to send email", "to", sub.Email, "error", err)
				continue
			}

			err = s.repository.UpdateLastTag(ctx, sub.ID, newTag)
			if err != nil {
				slog.Error("Scanner: Failed to update subscription tag", "sub_id", sub.ID, "error", err)
			}
		}
	}

}
