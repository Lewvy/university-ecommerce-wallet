package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"ecommerce/internal/validator"
	"encoding/base64"
	"math/big"
	"strconv"
	"time"
)

const ScopeActivation = "activation"

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func (t *Token) ValidateVerificationToken(v *validator.ValidationError, tokenPlaintext string) error {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 6, "token", "must be 6 bytes long")
	i, err := strconv.Atoi(tokenPlaintext)
	v.Check(err == nil, "token", "must be an number")
	v.Check(i >= 100000, "token", "invalid token: must be greater than 100000")
	v.Check(i <= 999999, "token", "invalid token: must be less than 999999")
	if len(v.Errors) == 0 {
		return nil
	}
	return v
}

func (t *Token) SetVerificationToken(token *Token) error {
	panic("unimplememted")
}

func (t *Token) DeleteToken(scope string, userID int64) error {
	panic("unimplememted")

}

func generateRandomString(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateAccessToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	result, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	token := &Token{
		Plaintext: result,
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	token.Hash = generateTokenHash(token.Plaintext)
	return token, nil
}

func MatchToken(token string, tokenHash []byte) bool {
	userTokenHash := generateTokenHash(token)
	return subtle.ConstantTimeCompare(userTokenHash, tokenHash) == 1
}

func GenerateVerificationToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	min := big.NewInt(100000)
	max := big.NewInt(999999)

	diff := new(big.Int).Sub(max, min)

	nBig, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return nil, err
	}
	result := new(big.Int).Add(nBig, min)

	token := &Token{
		Plaintext: result.String(),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}
	token.Hash = generateTokenHash(token.Plaintext)
	return token, nil
}

func generateTokenHash(tokenPlaintext string) []byte {

	hash := sha256.Sum256([]byte(tokenPlaintext))
	return hash[:]
}
