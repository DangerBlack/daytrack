package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const FAKE_OPAQUE_SLEEP_MILLISECONDS = 500

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[int(n.Int64())]
	}

	return string(b)
}

func SecureRandom(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))

	if err != nil {
		panic(err)
	}

	return nBig.Int64()
}

func FakeOpaqueOperation() {
	var err error
	var perturbation *big.Int
	var max *big.Int = big.NewInt(FAKE_OPAQUE_SLEEP_MILLISECONDS)

	if perturbation, err = rand.Int(rand.Reader, max); err != nil {
		perturbation = max
	}

	time.Sleep((200.0 + time.Duration(perturbation.Int64())) * time.Millisecond)
}

// PEMDecodePublicKey decodes an ed25519 PEM encoded with PKIX standard
func PEMDecodePublicKey(pub []byte) (ed25519.PublicKey, error) {
	pubBlock, _ := pem.Decode(pub)
	if pubBlock == nil {
		return nil, fmt.Errorf("no pem block found on public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}

	edPubKey, ok := pubKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an ed25519 key")
	}

	return edPubKey, nil
}

// PEMDecodePrivateKey decodes an ed25519 PEM encoded with PKCS8 standard
func PEMDecodePrivateKey(priv []byte) (ed25519.PrivateKey, error) {
	privBlock, _ := pem.Decode(priv)
	if privBlock == nil {
		return nil, fmt.Errorf("no pem block found on private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, err
	}

	edPrivKey, ok := privKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an ed25519 key")
	}

	return edPrivKey, nil
}

func PEMDecodeKeyPair(pub []byte, priv []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	edPubKey, err := PEMDecodePublicKey(pub)
	if err != nil {
		return nil, nil, err
	}

	edPrivKey, err := PEMDecodePrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	return edPubKey, edPrivKey, nil
}

func GenerateKeyPairFromSeed(seed []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	keypair := ed25519.NewKeyFromSeed(seed)

	pubKey, ok := keypair.Public().(ed25519.PublicKey)

	if !ok {
		return nil, nil, errors.New("public key is not of type ed25519")
	}

	return pubKey, keypair, nil
}
