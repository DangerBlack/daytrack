package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"512b.it/daytrack/src/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvironmentType string

const (
	Development EnvironmentType = "development"
	Production  EnvironmentType = "production"
)

type Configuration struct {
	Environment EnvironmentType
	HTTPHost    string
	HTTPPort    int

	Opaque OpaqueConfig
	JWT    JWTConfig
	Mail   Mail
}

type OpaqueConfig struct {
	SaltNonce      string
	ChallengeRange int
}

type JWTConfig struct {
	PublicKey  []byte
	PrivateKey []byte

	AccessTokenDuration time.Duration
}

type Mail struct {
	From     string
	Password string
	SMTP     string
	PORT     int
}

func NewConfiguration() Configuration {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	accessTokenDuration, err := time.ParseDuration(stringOrPanic("JWT_ACCESS_TOKEN_DURATION"))
	if err != nil {
		panic(err)
	}

	var jwtPubKey []byte
	var jwtPrivKey []byte

	encodedJWTPubKey := stringOrPanic("JWT_PUBLIC_KEY")
	encodedJWTPrivKey := stringOrPanic("JWT_PRIVATE_KEY")

	encodedJWTPubKey = strings.Replace(encodedJWTPubKey, `\n`, "\n", -1)
	encodedJWTPrivKey = strings.Replace(encodedJWTPrivKey, `\n`, "\n", -1)

	if encodedJWTPubKey != "" && encodedJWTPrivKey != "" {
		jwtPubKey, jwtPrivKey, err = utils.PEMDecodeKeyPair([]byte(encodedJWTPubKey), []byte(encodedJWTPrivKey))
	} else {
		log.Warn().Msg("Unable to read the keypair form env, generating new one")
		panic(fmt.Errorf("unable to read the keypair form env, generating new one"))
	}

	if err != nil {
		panic(err)
	}

	return Configuration{
		Environment: Development,
		HTTPHost:    stringOrPanic("HTTP_HOST"),
		HTTPPort:    intOrPanic("HTTP_PORT"),
		Opaque: OpaqueConfig{
			SaltNonce:      stringOrPanic("OPAQUE_SALT_NONCE"),
			ChallengeRange: intOrPanic("OPAQUE_CHALLENGE_RANGE"),
		},
		JWT: JWTConfig{
			PublicKey:           jwtPubKey,
			PrivateKey:          jwtPrivKey,
			AccessTokenDuration: accessTokenDuration,
		},
		Mail: Mail{
			From:     stringOrPanic("MAIL_FROM"),
			Password: stringOrPanic("MAIL_PASSWORD"),
			SMTP:     stringOrPanic("MAIL_SMTP_HOST"),
			PORT:     intOrPanic("MAIL_SMTP_PORT"),
		},
	}
}

func stringOrPanic(key string) string {
	var result, found = os.LookupEnv(key)

	if !found {
		panic(errors.New("configuration value not set for key: " + key))
	}

	return result
}

func intOrPanic(key string) int {
	var result, found = os.LookupEnv(key)

	if !found {
		panic(errors.New("configuration value not set for key: " + key))
	}

	intResult, err := strconv.ParseUint(result, 10, 32)
	if err != nil {
		panic(errors.New("configuration value for key: " + key + " is not a int"))
	}

	return int(intResult)
}
