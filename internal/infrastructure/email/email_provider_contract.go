package email

import "context"

type EmailProvider interface {
	SendReleaseNotification(ctx context.Context, toEmail, repository, tag string) error
}
