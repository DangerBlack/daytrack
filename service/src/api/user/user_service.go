package user

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"time"

	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
)

var ErrorOpaqueChallengeVerificationFailed error = errors.New("opaque challenge verification failed")

type Service struct {
	db            *database.Database
	configuration models.Configuration
	logger        utils.ContextLogger
}

func New(db *database.Database, configuration models.Configuration) *Service {
	logger := utils.InitServiceLogger("UserService")

	return &Service{
		db:            db,
		configuration: configuration,
		logger:        logger,
	}
}

func (s *Service) GenerateSalt(ctx context.Context, saltNonce, email string) string {
	h := sha256.New()
	h.Write([]byte(email + saltNonce))

	return b64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *Service) CreateUser(ctx context.Context, username, email, password string) (int64, error) {
	salt := s.GenerateSalt(context.Background(), s.configuration.Opaque.SaltNonce, email)

	id, err := s.db.InsertUser(username, email, password, salt)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) SigninUser(ctx context.Context, email, challenge, signedChallengeBase64 string) (*models.Token, error) {
	var err error
	var user *models.User
	var signedChallenge []byte
	var publicKey []byte

	if user, err = s.db.GetUserByEmail(email); err != nil {
		return nil, err
	}

	if signedChallenge, err = base64.StdEncoding.DecodeString(signedChallengeBase64); err != nil {
		utils.FakeOpaqueOperation()
		return nil, err
	}

	if publicKey, err = b64.StdEncoding.DecodeString(user.PublicKey); err != nil {
		utils.FakeOpaqueOperation()
		return nil, err
	}

	if !ed25519.Verify(ed25519.PublicKey(publicKey), []byte(challenge), []byte(signedChallenge)) {
		return nil, ErrorOpaqueChallengeVerificationFailed
	}

	var token string
	if token, err = models.GenerateAccessJWT(s.configuration.JWT.PrivateKey, time.Duration(s.configuration.JWT.AccessTokenDuration), fmt.Sprintf("%d", user.ID)); err != nil {
		utils.FakeOpaqueOperation()
		return nil, err
	}

	return &models.Token{
		Token:   token,
		ExpDate: time.Now().Add(time.Duration(s.configuration.JWT.AccessTokenDuration)),
	}, nil
}

func (s *Service) SendMagicLink(ctx context.Context, email string) error {
	var err error
	var key *string
	var userID int64
	if _, err = s.CreateUser(ctx, "placeholder", email, "~~~"); err != nil {
		if !errors.Is(err, database.ErrorDuplicate) {
			return err
		}
	}

	name := utils.GenerateRandomString(32)

	if key, err = s.db.InsertAPIKey(userID, name); err != nil {
		return err
	}

	url := "https://dailytrack.io"
	if s.configuration.Environment == "development" {
		url = "http://10.0.2.2:3000"
	}

	url = fmt.Sprintf("%s/v1/users/magic-link?key=%s", url, *key)

	if err = models.SendMagicLink(ctx, s.configuration, email, url); err != nil {
		return err
	}

	return nil
}
