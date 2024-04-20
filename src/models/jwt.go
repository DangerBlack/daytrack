package models

import (
	"context"
	"crypto/ed25519"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTPubKeyType = ed25519.PublicKey
type JWTPrivKeyType = ed25519.PrivateKey

var ErrWrongTokenSigning error = errors.New("wrong token signing method")
var ErrTokenNotValid error = errors.New("token not valid")

const (
	JWTExpClaimKey     = "exp" // Expiration
	JWTSubjectClaimKey = "sub" // AccountID
)

const USER_ID_CONTEXT_KEY = "UserID"

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

func ParseAndValidateJWT(ctx context.Context, stringToken string, key JWTPubKeyType) (*jwt.Token, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodEd25519)
		if !ok {
			return nil, ErrWrongTokenSigning
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenNotValid
	}

	return token, nil
}
