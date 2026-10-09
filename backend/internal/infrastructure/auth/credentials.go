package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"traci/backend/internal/application"
)

type Credentials struct {
	secret   []byte
	lifetime time.Duration
}

var _ application.Credentials = (*Credentials)(nil)

func NewCredentials(secret string, lifetime time.Duration) (*Credentials, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must contain at least 32 bytes")
	}
	if lifetime < time.Second {
		return nil, fmt.Errorf("TOKEN_TTL must be at least one second")
	}
	return &Credentials{secret: []byte(secret), lifetime: lifetime}, nil
}

func (c *Credentials) Hash(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, 600000, 32)
	if err != nil {
		return "", err
	}
	return "pbkdf2-sha256$600000$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(key), nil
}

func (c *Credentials) Verify(password, hash string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2-sha256" || parts[1] != "600000" {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil || len(salt) != 16 {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil || len(expected) != 32 {
		return false
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, 600000, 32)
	return err == nil && subtle.ConstantTimeCompare(key, expected) == 1
}

func (c *Credentials) Issue(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject: userID.String(), Issuer: "traci", Audience: jwt.ClaimStrings{"traci"},
		IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(c.lifetime)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(c.secret)
}

func (c *Credentials) Parse(raw string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		return c.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired(), jwt.WithIssuer("traci"), jwt.WithAudience("traci"), jwt.WithIssuedAt())
	if err != nil || !token.Valid {
		return uuid.Nil, application.ErrUnauthorized
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, application.ErrUnauthorized
	}
	return id, nil
}
