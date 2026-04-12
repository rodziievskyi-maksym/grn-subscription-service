package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/config"
)

type SmtpProvider struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSmtpProvider() *SmtpProvider {
	return &SmtpProvider{
		host:     config.Cfg().SmtpHost,
		port:     config.Cfg().SmtpPort,
		username: config.Cfg().SmtpUser,
		password: config.Cfg().SmtpPass,
		from:     config.Cfg().SmtpFrom,
	}
}

func (p *SmtpProvider) SendReleaseNotification(ctx context.Context, toEmail, repository, tag string) error {
	subject := fmt.Sprintf("Subject: New Release for %s!\n", repository)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>New version detected!</h2>
				<p>The repository <strong>%s</strong> has a new release: <strong>%s</strong></p>
				<br>
				<p><a href="https://github.com/%s/releases/tag/%s">View Release on GitHub</a></p>
			</body>
		</html>`, repository, tag, repository, tag)

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%s", p.host, p.port)
	auth := smtp.PlainAuth("", p.username, p.password, p.host)

	// Відправка
	err := smtp.SendMail(addr, auth, p.from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("smtp send error: %w", err)
	}

	return nil
}
