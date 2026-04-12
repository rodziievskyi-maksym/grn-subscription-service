package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Repository  string    `json:"repository"`
	LastSeenTag string    `json:"last_seen_tag"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewSubscription(email, repoOwner, repoName, lastSeenTag string) *Subscription {
	return &Subscription{
		ID:          uuid.New(),
		Email:       email,
		Repository:  JoinRepoOwnerAndName(repoOwner, repoName),
		LastSeenTag: lastSeenTag,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func JoinRepoOwnerAndName(repoOwner, repoName string) string {
	return fmt.Sprintf("%s/%s", repoOwner, repoName)
}
