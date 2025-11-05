package cache

import (
	"context"
	"time"
)

type Cache interface {
	SetVerificationToken(ctx context.Context, tokenHash string, userID int64, expiry time.Duration) error
	AddEmailToQueue(ctx context.Context, email, data string) error
	GetUserIDByToken(ctx context.Context, tokenHash string) (int64, error)
	DeleteToken(ctx context.Context, tokenHash string) error
}
