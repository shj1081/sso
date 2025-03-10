package service

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/shj1081/sso/internal/config"
	"github.com/shj1081/sso/internal/storer"
)

// EmailService 구조체 정의
type EmailService struct {
	cfg *config.Config
	st  storer.Storer
}

// NewEmailService 생성자 함수
func NewEmailService(cfg *config.Config, st storer.Storer) *EmailService {
	return &EmailService{
		cfg: cfg,
		st:  st,
	}
}

func (es *EmailService) GetVerifyCodeByID(ctx context.Context, userId int64) (string, error) {
	return es.st.GetVerifyCodeByID(ctx, userId)
}

// 이메일 발송 함수
func (es *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", es.cfg.SMTPUser, es.cfg.SMTPPassword, es.cfg.SMTPHost)

	// 이메일 메시지 구성
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	recipients := strings.Split(to, ",")

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", es.cfg.SMTPHost, es.cfg.SMTPPort),
		auth,
		es.cfg.SMTPFrom,
		recipients,
		[]byte(msg),
	)
	return err
}

// 인증 코드 저장 및 이메일 전송
func (es *EmailService) SendVerificationEmail(ctx context.Context, email string, code string) error {

	// 이메일 전송
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", code)
	if err := es.SendEmail(email, subject, body); err != nil {
		return fmt.Errorf("failed to send email")
	}

	return nil
}

// 인증 코드 검증
func (es *EmailService) VerifyCode(ctx context.Context, code string, ans string, userId int64) error {
	if code == ans {
		// userType skkuin으로 변경
		_, err := es.st.UpdateUser(ctx, &storer.User{
			ID:       userId,
			UserType: "skkuin",
		})
		if err != nil {
			return fmt.Errorf("failed to update usertype")
		}

		return nil
	} else {
		return fmt.Errorf("invalid code")
	}

}
