package email

import "context"

type Sender interface {
	SendVerificationEmail(ctx context.Context, toEmail, toName, token string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, toName, token string) error
	SendReportReadyEmail(ctx context.Context, toEmail, toName, carName, downloadURL string) error
}
