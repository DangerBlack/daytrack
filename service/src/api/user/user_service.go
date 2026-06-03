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
	var userID string
	valid := true
	found := true

	user, err = s.db.GetUserByEmail(email)
	if err != nil {
		found = false
		userID = "0"
	} else {
		userID = fmt.Sprintf("%d", user.ID)
	}

	if signedChallenge, err = base64.StdEncoding.DecodeString(signedChallengeBase64); err != nil {
		return nil, err
	}

	if !found {
		fakeSeed := sha256.Sum256([]byte(email + "__timing_mask__"))
		fakeKey := ed25519.NewKeyFromSeed(fakeSeed[:])
		ed25519.Verify(fakeKey.Public().(ed25519.PublicKey), []byte(challenge), signedChallenge)
		return nil, ErrorOpaqueChallengeVerificationFailed
	}

	if publicKey, err = b64.StdEncoding.DecodeString(user.PublicKey); err != nil {
		return nil, err
	}

	if !ed25519.Verify(ed25519.PublicKey(publicKey), []byte(challenge), signedChallenge) {
		valid = false
	}

	var token string
	if token, err = models.GenerateAccessJWT(s.configuration.JWT.PrivateKey, time.Duration(s.configuration.JWT.AccessTokenDuration), userID); err != nil {
		return nil, err
	}

	if !valid {
		return nil, ErrorOpaqueChallengeVerificationFailed
	}

	return &models.Token{
		Token:    token,
		ExpDate:  time.Now().Add(time.Duration(s.configuration.JWT.AccessTokenDuration)),
		Username: user.Username,
	}, nil
}
