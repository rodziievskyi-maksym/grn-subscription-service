package email

type Provider interface {
	SendReleaseNotification(toEmail, repository, tag string) error
}
