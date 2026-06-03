package models

import (
	"crypto/ed25519"
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
	DBPath      string

	Opaque              OpaqueConfig
	JWT                 JWTConfig
	RateLimitBurst      int
	RateLimitInterval   time.Duration
	SignupDisabled      bool
}

type OpaqueConfig struct {
	SaltNonce      string
	ChallengeRange int
}

type JWTConfig struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey

	AccessTokenDuration time.Duration
}

func NewConfiguration() Configuration {
	_ = godotenv.Load()

	env := getEnv("ENVIRONMENT", "development")
	envType := Development
	if env == "production" {
		envType = Production
	}

	accessTokenDuration, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"))
	if err != nil {
		accessTokenDuration = 15 * time.Minute
	}

	var jwtPubKey ed25519.PublicKey
	var jwtPrivKey ed25519.PrivateKey

	encodedJWTPubKey := getEnv("JWT_PUBLIC_KEY", "")
	encodedJWTPrivKey := getEnv("JWT_PRIVATE_KEY", "")

	encodedJWTPubKey = strings.Replace(encodedJWTPubKey, `\n`, "\n", -1)
	encodedJWTPrivKey = strings.Replace(encodedJWTPrivKey, `\n`, "\n", -1)

	if encodedJWTPubKey != "" && encodedJWTPrivKey != "" {
		var pubKeyBytes, privKeyBytes []byte
		pubKeyBytes, privKeyBytes, err = utils.PEMDecodeKeyPair([]byte(encodedJWTPubKey), []byte(encodedJWTPrivKey))
		if err != nil {
			log.Warn().Err(err).Msg("Failed to decode JWT key pair, generating ephemeral keys")
		} else {
			jwtPubKey = pubKeyBytes
			jwtPrivKey = privKeyBytes
		}
	}

	if jwtPrivKey == nil {
		log.Warn().Msg("No valid JWT key pair found, generating ephemeral keys")
		_, jwtPrivKey, err = ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}
		jwtPubKey = jwtPrivKey.Public().(ed25519.PublicKey)
	}

	rateLimitInterval, err := time.ParseDuration(getEnv("RATE_LIMIT_INTERVAL", "1s"))
	if err != nil {
		rateLimitInterval = time.Second
	}

	return Configuration{
		Environment: envType,
		HTTPHost:    getEnv("HTTP_HOST", "localhost"),
		HTTPPort:    getEnvInt("HTTP_PORT", 3000),
		DBPath:      getEnv("DB_PATH", "./archive/database.db"),
		Opaque: OpaqueConfig{
			SaltNonce:      getEnv("OPAQUE_SALT_NONCE", "default-salt-nonce"),
			ChallengeRange: getEnvInt("OPAQUE_CHALLENGE_RANGE", 100000),
		},
		JWT: JWTConfig{
			PublicKey:           jwtPubKey,
			PrivateKey:          jwtPrivKey,
			AccessTokenDuration: accessTokenDuration,
		},
		RateLimitBurst:    getEnvInt("RATE_LIMIT_BURST", 3),
		RateLimitInterval: rateLimitInterval,
		SignupDisabled:    getEnv("DISABLE_SIGNUP", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return intVal
}


