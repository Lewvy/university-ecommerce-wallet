package cache

import (
	"context"
	"time"
)

type Cache interface {
	SetVerificationToken(ctx context.Context, tokenHash string, userID int64, expiry time.Duration) error
	AddEmailToQueue(ctx context.Context, email, data string) error
	GetUserIDByTokenHash(ctx context.Context, tokenHash string) (int64, error)
	DeleteToken(ctx context.Context, tokenHash string) error
	GetTokenHashByUserID(ctx context.Context, userID int64) (string, error)
}
