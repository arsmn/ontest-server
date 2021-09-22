package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/settings"
	"github.com/arsmn/ontest-server/user"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	ResetPasswordTmpl     = "reset_password.html"
	EmailVerificationTmpl = "email_verification.html"

	ResetPasswordSubject     = "Reset Password"
	EmailVerificationSubject = "Verify Email"
)

type (
	smtpDependencies interface {
		settings.Provider
		xlog.Provider
	}
	SMTP struct {
		auth smtp.Auth
		dx   smtpDependencies
		t    *template.Template
	}
)

func NewMailerSMTP(dx smtpDependencies) (*SMTP, error) {
	c := dx.Settings().Mail.SMTP

	t, err := template.ParseFiles(
		"./module/mail/reset_password.html",
		"./module/mail/email_verification.html",
	)
	if err != nil {
		return nil, err
	}

	a := smtp.PlainAuth("", c.From, c.Password, c.Host)
	return &SMTP{a, dx, t}, nil
}

func (s *SMTP) sendMail(message []byte, to []string) error {
	c := s.dx.Settings().Mail.SMTP

	addr := c.Host + ":" + c.Port
	return smtp.SendMail(addr, s.auth, c.From, to, message)
}

func (s *SMTP) sendTemplate(tmpl string, data interface{}, sbj string, to []string) error {
	buffer := new(bytes.Buffer)
	buffer.WriteString("To: " + to[0] + "\r\nSubject: " + sbj + "\r\n" + MIME + "\r\n")
	if err := s.t.ExecuteTemplate(buffer, tmpl, data); err != nil {
		return err
	}

	return s.sendMail(buffer.Bytes(), to)
}

func (s *SMTP) SendResetPassword(ctx context.Context, u *user.User, code string) {
	go func() {
		err := s.sendTemplate(ResetPasswordTmpl, map[string]string{
			"URL":  fmt.Sprintf("%s/reset-password/%s", s.dx.Settings().Client.WebURL, code),
			"Name": u.FirstName,
		}, ResetPasswordSubject, []string{u.Email})

		if err != nil {
			s.dx.Logger().Warn(fmt.Sprintf("Error while sending reset password email for: %s", u.Email), xlog.Err(err))
		}
	}()
}

func (s *SMTP) SendVerification(ctx context.Context, u *user.User, code string) {
	go func() {
		err := s.sendTemplate(EmailVerificationTmpl, map[string]string{
			"URL":  fmt.Sprintf("%s/verify-email/%s", s.dx.Settings().Client.WebURL, code),
			"Name": u.FirstName,
		}, EmailVerificationSubject, []string{u.Email})

		if err != nil {
			s.dx.Logger().Warn(fmt.Sprintf("Error while sending email verification email for: %s", u.Email), xlog.Err(err))
		}
	}()
}
