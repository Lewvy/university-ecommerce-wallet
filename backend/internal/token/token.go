package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"ecommerce/internal/validator"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("Error: jwt environment variable not set.")
	}
	return secret
}

var ErrTokenNotFound = errors.New("token not found")
var ErrInvalidToken = errors.New("invalid or expired token")

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopeRefresh        = "refresh"
)

type Token struct {
	Plaintext string
	Hash      string
	UserID    int64
	Expiry    time.Time
	Scope     string
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Scope  string `json:"scope"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	expirationTime := time.Now().Add(ttl)

	claims := &Claims{
		UserID: userID,
		Scope:  scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tkn.SignedString([]byte(JWTSecret()))
	if err != nil {
		return nil, err
	}

	return &Token{
		Plaintext: tokenString,
		UserID:    userID,
		Expiry:    expirationTime,
		Scope:     scope,
	}, nil
}

func GenerateRefreshToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
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

	token.Hash = GenerateTokenHash(token.Plaintext)
	return token, nil
}

func GenerateVerificationToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, err
	}
	result := nBig.Int64() + 100000

	token := &Token{
		Plaintext: strconv.FormatInt(result, 10),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}
	token.Hash = GenerateTokenHash(token.Plaintext)
	return token, nil
}

func VerifyAccessToken(tokenPlaintext string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(JWTSecret()), nil
	}

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenPlaintext, claims, keyFunc)
	if err != nil || !tkn.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Scope != ScopeAuthentication {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func MatchToken(token string, storedHash string) (bool, error) {
	userTokenHashStr := GenerateTokenHash(token)

	userHashBytes, err := hex.DecodeString(userTokenHashStr)
	if err != nil {
		return false, err
	}

	storedHashBytes, err := hex.DecodeString(storedHash)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(userHashBytes, storedHashBytes) == 1, nil
}

func GenerateTokenHash(tokenPlaintext string) string {
	hash := sha256.Sum256([]byte(tokenPlaintext))
	return hex.EncodeToString(hash[:])
}

func generateRandomString(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
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
