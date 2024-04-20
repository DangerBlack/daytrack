package client

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"512b.it/daytrack/src/models"
	user_utils "512b.it/daytrack/src/utils"
	"512b.it/daytrack/tests/utils"
)

const BASE_URL = "http://localhost:3000"

func GenerateChallenge(email string) (*models.Challenge, error) {
	var err error
	var response models.Challenge
	url := BASE_URL + "/v1/user/challenge?email=" + email

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodGet),
		utils.WithExpectedStatusCode(http.StatusOK),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create get account request: %w", err)
	}

	return &response, nil
}

func SignUp(username, email, password string) error {
	var err error
	var challenge *models.Challenge
	url := BASE_URL + "/v1/user/signup"

	if challenge, err = GenerateChallenge(email); err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write([]byte(password + challenge.Salt))
	seed := hash.Sum(nil)

	var publicKey ed25519.PublicKey
	if publicKey, _, err = user_utils.GenerateKeyPairFromSeed(seed); err != nil {
		return err
	}

	print("public key", base64.StdEncoding.EncodeToString(publicKey))
	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithRequestBody(map[string]string{
			"username":   username,
			"email":      email,
			"public_key": base64.StdEncoding.EncodeToString(publicKey),
			"salt":       challenge.Salt,
		}),
		utils.WithExpectedStatusCode(http.StatusCreated),
	); err != nil {
		return fmt.Errorf("failed unable to create an user request: %w", err)
	}

	return nil
}

func SignIn(email, password string) (*models.Token, error) {
	var err error
	var response models.Token
	var challenge *models.Challenge
	url := BASE_URL + "/v1/user/signin"

	if challenge, err = GenerateChallenge(email); err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write([]byte(password + challenge.Salt))
	seed := hash.Sum(nil)

	var privateKey ed25519.PrivateKey
	if _, privateKey, err = user_utils.GenerateKeyPairFromSeed(seed); err != nil {
		return nil, err
	}

	signedChallenge := ed25519.Sign(privateKey, []byte(challenge.Challenge))

	if err = utils.DoRequest(
		url,
		utils.WithRequestMethod(http.MethodPost),
		utils.WithRequestBody(map[string]string{
			"email":            email,
			"challenge":        challenge.Challenge,
			"signed_challenge": base64.StdEncoding.EncodeToString(signedChallenge),
		}),
		utils.WithExpectedStatusCode(http.StatusOK),
		utils.ExtractGenericModel(&response),
	); err != nil {
		return nil, fmt.Errorf("failed unable to create an user request: %w", err)
	}

	return &response, nil
}
