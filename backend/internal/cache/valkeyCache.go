package cache

import (
	"context"
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

const (
	User         = "user"
	Verification = "verification"
)

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

func (v *ValkeyCache) GetTokenHashByUserID(ctx context.Context, userID int64) (string, error) {

	key := fmt.Sprintf("%s:%d:%s", User, userID, Verification)
	tokenHash, err := v.Client.Do(ctx, v.Client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return "", err
	}
	return tokenHash, nil
}

func (v *ValkeyCache) SetVerificationToken(ctx context.Context, tokenHash string, userID int64, expiry time.Duration) error {
	userIDstr := strconv.FormatInt(userID, 10)

	token := fmt.Sprintf("%s:%s:%s", "verification", "hash", tokenHash)

	tokenKeyCmd := v.Client.B().Set().Key(token).Value(userIDstr).Nx().Ex(expiry).Build()

	userVerificationKey := fmt.Sprintf("user:%s:%s", userIDstr, "verification")
	userKeyCmd := v.Client.B().Set().Key(userVerificationKey).Value(tokenHash).Ex(expiry).Build()

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
func (v *ValkeyCache) cartKey(userID int64) string {
	return fmt.Sprintf("cart:%d", userID)
}

func (v *ValkeyCache) AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	key := v.cartKey(userID)
	field := strconv.FormatInt(productID, 10)

	cmd := v.Client.B().Hincrby().Key(key).Field(field).Increment(int64(quantity)).Build()
	return v.Client.Do(ctx, cmd).Error()
}

func (v *ValkeyCache) GetCart(ctx context.Context, userID int64) (map[string]string, error) {
	key := v.cartKey(userID)
	cmd := v.Client.B().Hgetall().Key(key).Build()

	res, err := v.Client.Do(ctx, cmd).AsStrMap()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	return res, nil
}

func (v *ValkeyCache) UpdateCartItemQuantity(ctx context.Context, userID int64, productID int64, quantity int) error {
	key := v.cartKey(userID)
	field := strconv.FormatInt(productID, 10)

	if quantity <= 0 {
		return v.DeleteCartItem(ctx, userID, productID)
	}

	cmd := v.Client.B().Hset().Key(key).FieldValue().FieldValue(field, strconv.Itoa(quantity)).Build()
	return v.Client.Do(ctx, cmd).Error()
}

func (v *ValkeyCache) DeleteCartItem(ctx context.Context, userID int64, productID int64) error {
	key := v.cartKey(userID)
	field := strconv.FormatInt(productID, 10)

	cmd := v.Client.B().Hdel().Key(key).Field().Field(field).Build()
	return v.Client.Do(ctx, cmd).Error()
}

func (v *ValkeyCache) ClearCart(ctx context.Context, userID int64) error {
	key := v.cartKey(userID)
	cmd := v.Client.B().Del().Key(key).Build()
	return v.Client.Do(ctx, cmd).Error()
}

func (v *ValkeyCache) GetCartCount(ctx context.Context, userID int64) (int64, error) {
	key := v.cartKey(userID)
	cmd := v.Client.B().Hlen().Key(key).Build()
	return v.Client.Do(ctx, cmd).AsInt64()
}
