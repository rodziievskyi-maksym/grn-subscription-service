//go:build integration

package github_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/config"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
)

const (
	GitHubTestOwner = "gin-gonic"
	GitHubTestRepo  = "gin"
)

var testToken string

func TestMain(m *testing.M) {
	v := validator.New()

	if err := config.NewConfig(v, "../../../.env"); err != nil {
		log.Fatal(err)
	}

	testToken = config.Cfg().GitHubToken

	os.Exit(m.Run())
}

func TestGetLatestTag(t *testing.T) {
	client := github.NewClient(testToken)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tag, err := client.GetLatestTag(ctx, GitHubTestOwner, GitHubTestRepo)
	if err != nil {
		t.Fatalf("Failed to get latest tag: %v", err)
	}

	t.Logf("Latest tag: %s", tag)
}
