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
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", code)
	if err := es.SendEmail(email, subject, body); err != nil {
		return fmt.Errorf("failed to send email")
	}
	return nil
}

// 세션 기반 인증 코드 전송
func (es *EmailService) SendVerificationEmailBySession(ctx context.Context, sessionID string, email string) error {
	// 세션에서 데이터 가져오기
	sd, err := es.st.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session db error: %w", err)
	}
	if sd == nil {
		return fmt.Errorf("invalid session")
	}

	return es.SendVerificationEmail(ctx, email, sd.VerifyCode)
}

// userID 기반 인증 코드 전송
func (es *EmailService) SendVerificationEmailByUserID(ctx context.Context, userID int64, email string) error {
	code, err := es.st.GetVerifyCodeByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get verification code: %w", err)
	}

	return es.SendVerificationEmail(ctx, email, code)
}

// 인증 코드 검증 및 유저 업데이트
func (es *EmailService) VerifyCode(ctx context.Context, code string, ans string, userId int64) error {
	if code == ans {
		_, err := es.st.UpdateUser(ctx, &storer.User{
			ID:       userId,
			UserType: "skkuin",
		})
		if err != nil {
			return fmt.Errorf("failed to update usertype")
		}
		return nil
	}
	return fmt.Errorf("invalid code")
}

// 세션 기반 인증 코드 검증
func (es *EmailService) VerifyCodeBySession(ctx context.Context, sessionID, code string) error {
	sd, err := es.st.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session db error: %w", err)
	}
	if sd == nil {
		return fmt.Errorf("invalid session")
	}

	return es.VerifyCode(ctx, code, sd.VerifyCode, sd.UserId)
}

// userID 기반 인증 코드 검증
func (es *EmailService) VerifyCodeByUserID(ctx context.Context, userID int64, code string) error {
	verifyCode, err := es.st.GetVerifyCodeByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get verification code: %w", err)
	}

	return es.VerifyCode(ctx, code, verifyCode, userID)
}
