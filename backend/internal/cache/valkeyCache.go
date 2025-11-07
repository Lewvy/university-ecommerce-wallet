package cache

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"
)

type ValkeyCache struct {
	Client valkey.Client
}

func (v *ValkeyCache) SetVerificationToken(ctx context.Context, tokenHash []byte, userID int64, expiry time.Duration) error {
	userIDstr := strconv.FormatInt(userID, 10)

	encodedHash := base64.RawStdEncoding.EncodeToString(tokenHash)

	token := fmt.Sprintf("%s:%s:%s", "verification", "hash", encodedHash)

	tokenKeyCmd := v.Client.B().Set().Key(token).Value(userIDstr).Nx().Ex(expiry).Build()

	userVerificationKey := fmt.Sprintf("user:%s:%s", userIDstr, "verification")
	userKeyCmd := v.Client.B().Set().Key(userVerificationKey).Value(encodedHash).Ex(expiry).Build()

	resp := v.Client.DoMulti(ctx, tokenKeyCmd, userKeyCmd)
	if err := resp[0].Error(); err != nil && !valkey.IsValkeyNil(err) {
		return fmt.Errorf("valkey transaction failed on SET command: %w", err)
	}

	if valkey.IsValkeyNil(resp[0].Error()) {
		cleanupCmd := v.Client.B().Del().Key(userVerificationKey).Build()
		if cleanupErr := v.Client.Do(ctx, cleanupCmd).Error(); cleanupErr != nil {
			return fmt.Errorf("token hash collision (NX failed), AND cleanup failed: %w, cleanup error: %v",
				errors.New("token hash collision detected"), cleanupErr)
		}
		return errors.New("token hash collision detected (NX condition failed)")
	}

	if err := resp[1].Error(); err != nil {
		return fmt.Errorf("valkey command 2 error on user key SET: %w", err)
	}

	return nil
}

const EmailQueueKey = "queue:emails"

func (v *ValkeyCache) AddEmailToQueue(ctx context.Context, email, jobJSON string) error {
	err := v.Client.Do(ctx,
		v.Client.B().Lpush().Key(EmailQueueKey).Element(jobJSON).Build(),
	).Error()
	return err
}

func (v *ValkeyCache) GetUserIDByTokenHash(ctx context.Context, tokenHash string) (int64, error) {
	id, err := v.Client.Do(ctx, v.Client.B().Get().Key(tokenHash).Build()).AsInt64()

	if err != nil {
		if valkey.IsValkeyNil(err) {
			return 0, errors.New("verification token not found or expired")
		}
		return 0, fmt.Errorf("failed to retrieve user ID for token: %w", err)
	}

	return id, nil
}
func (v *ValkeyCache) DeleteToken(ctx context.Context, tokenHash string) error {
	id, err := v.GetUserIDByTokenHash(ctx, tokenHash)

	if err != nil {
		if valkey.IsValkeyNil(err) || err.Error() == "verification token not found or expired" {
			return nil
		}
		return err
	}

	userIDstr := strconv.FormatInt(id, 10)

	tokenDelCmd := v.Client.B().Del().Key(tokenHash).Build()
	userDelCmd := v.Client.B().Del().Key(userIDstr).Build()

	resp := v.Client.DoMulti(ctx, tokenDelCmd, userDelCmd)

	if err := resp[0].Error(); err != nil && !valkey.IsValkeyNil(err) {
		return fmt.Errorf("valkey transaction failed during primary DEL: %w", err)
	}

	if err := resp[1].Error(); err != nil && !valkey.IsValkeyNil(err) {
		return fmt.Errorf("valkey transaction failed during secondary DEL: %w", err)
	}

	return nil
}

func NewValkeyCache() (*ValkeyCache, error) {
	address, err := valkey.ParseURL(os.Getenv("CACHE_DSN"))
	if err != nil {
		return nil, err
	}
	client, err := valkey.NewClient(address)
	if err != nil {
		return nil, err
	}
	return &ValkeyCache{Client: client}, nil
}
