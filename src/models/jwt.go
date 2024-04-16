package models

import (
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTPubKeyType = ed25519.PublicKey
type JWTPrivKeyType = ed25519.PrivateKey

const (
	JWTExpClaimKey     = "exp"  // Expiration
	JWTSubjectClaimKey = "sub"  // AccountID
	JWTTypeClaimKey    = "type" // one of JWTTokenType
)

func GenerateAccessJWT(
	key JWTPrivKeyType,
	exp time.Duration,
	subjectID string,
) (string, error) {
	return generateJWT(
		key,
		withJWTExpirationClaim(exp),
		withJWTSubjectIDClaim(subjectID),
	)
}

func generateJWT(
	key JWTPrivKeyType,
	claimMappers ...func(jwt.MapClaims),
) (string, error) {
	var err error
	var tokenString string

	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)

	for _, mapper := range claimMappers {
		mapper(claims)
	}

	if tokenString, err = token.SignedString(key); err != nil {
		return "", err
	}

	return tokenString, nil
}

func withJWTClaim(key string, value any) func(jwt.MapClaims) {
	return func(mc jwt.MapClaims) {
		mc[key] = value
	}
}

func withJWTExpirationClaim(exp time.Duration) func(jwt.MapClaims) {
	return withJWTClaim(JWTExpClaimKey, time.Now().Add(exp).Unix())
}

func withJWTSubjectIDClaim(accountID string) func(jwt.MapClaims) {
	return withJWTClaim(JWTSubjectClaimKey, accountID)
}
