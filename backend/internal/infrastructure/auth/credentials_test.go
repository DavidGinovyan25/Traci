package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	secret := strings.Repeat("a", 32)
	credentials, err := NewCredentials(secret, time.Hour)
	require.NoError(t, err)
	password := strings.Repeat("p", 128)
	first, err := credentials.Hash(password)
	require.NoError(t, err)
	second, err := credentials.Hash(password)
	require.NoError(t, err)
	require.NotEqual(t, second, first, "password verification failed")
	assert.True(t, credentials.Verify(password, first), "password verification failed")
	assert.False(t, credentials.Verify("wrong", first), "password verification failed")
	assert.False(t, credentials.Verify(password, "invalid"), "password verification failed")
	id := uuid.New()
	token, err := credentials.Issue(id)
	require.NoError(t, err)
	parsed, err := credentials.Parse(token)
	require.NoError(t, err, "token round trip: %v", err)
	assert.Equal(t, id, parsed, "token round trip: %v", err)
	for _, claims := range []jwt.RegisteredClaims{
		{Subject: id.String(), Issuer: "traci", Audience: jwt.ClaimStrings{"traci"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour))},
		{Subject: id.String(), Issuer: "traci", Audience: jwt.ClaimStrings{"traci"}},
		{Subject: id.String(), Issuer: "wrong", Audience: jwt.ClaimStrings{"traci"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	} {
		raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		require.NoError(t, err)
		{
			_, err := credentials.Parse(raw)
			require.Error(t, err, "invalid claims accepted")
		}
	}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.RegisteredClaims{
		Subject: id.String(), Issuer: "traci", Audience: jwt.ClaimStrings{"traci"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte(secret))
	require.NoError(t, err)
	{
		_, err := credentials.Parse(raw)
		require.Error(t, err, "unexpected signing method accepted")
	}
}
