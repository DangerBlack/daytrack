package tests

import (
	"testing"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/tests/client"
	faker "512b.it/daytrack/tests/utils"
)

func TestChallenge(t *testing.T) {
	t.Parallel()

	var err error
	var challenge, challenge2 *models.Challenge

	email := faker.Faker.Email()

	if challenge, err = client.GenerateChallenge(email); err != nil {
		t.Fatalf("unable to generate a challenge %v", err)
	}

	if challenge.Challenge == "" {
		t.Fatal("Failed unable extract a valid challenge")
		return
	}

	if challenge.Salt == "" {
		t.Fatal("Failed unable extract a valid salt")
		return
	}

	if challenge2, err = client.GenerateChallenge(email); err != nil {
		t.Fatalf("unable to generate a challenge %v", err)
		return
	}

	if challenge.Challenge == challenge2.Challenge {
		t.Fatal("Failed challenge should not be equals")
		return
	}

	if challenge.Salt != challenge2.Salt {
		t.Fatal("Failed salt should be equals")
		return
	}
}

func TestSignUp(t *testing.T) {
	t.Parallel()

	username := faker.Faker.Username()
	email := faker.Faker.Email()
	password := faker.Faker.Password()

	if err := client.SignUp(username, email, password); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	var err error
	var token *models.Token
	username := faker.Faker.Username()
	email := faker.Faker.Email()
	password := faker.Faker.Password()

	if err = client.SignUp(username, email, password); err != nil {
		t.Fatalf("unable to sign up %v", err)
		return
	}

	if token, err = client.SignIn(email, password); err != nil {
		t.Fatalf("unable to sign in %v", err)
		return
	}

	if token == nil {
		t.Fatal("Failed unable to extract a valid token")
		return
	}

	if token.Token == "" {
		t.Fatal("Failed unable to extract a valid token")
		return
	}
}
