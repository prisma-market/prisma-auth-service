package email

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *EmailService) SendEmail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	headers := map[string]string{
		"From":         s.from,
		"To":           to[0],
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, to, []byte(message))
}

// 비밀번호 재설정 이메일 템플릿
func (s *EmailService) SendPasswordResetEmail(to string, resetToken string, resetURL string) error {
	subject := "비밀번호 재설정"
	body := fmt.Sprintf(`
        <h2>비밀번호 재설정</h2>
        <p>아래 링크를 클릭하여 비밀번호를 재설정하세요:</p>
        <p><a href="%s?token=%s">비밀번호 재설정</a></p>
        <p>이 링크는 1시간 동안 유효합니다.</p>
        <p>본인이 요청하지 않았다면 이 이메일을 무시하세요.</p>
    `, resetURL, resetToken)

	return s.SendEmail([]string{to}, subject, body)
}

// 이메일 인증 템플릿
func (s *EmailService) SendVerificationEmail(to string, verificationToken string, verifyURL string) error {
	subject := "이메일 주소 인증"
	body := fmt.Sprintf(`
        <h2>이메일 주소 인증</h2>
        <p>아래 링크를 클릭하여 이메일 주소를 인증하세요:</p>
        <p><a href="%s?token=%s">이메일 인증</a></p>
        <p>이 링크는 24시간 동안 유효합니다.</p>
    `, verifyURL, verificationToken)

	return s.SendEmail([]string{to}, subject, body)
}
