package jwt

import (
	"crypto/sha256"

	"github.com/google/uuid"
	"golang.org/x/crypto/ed25519"
)

func mustGenerateKID(publicKey ed25519.PublicKey) string {
	hash := sha256.New()
	_, err := hash.Write(publicKey)
	if err != nil {
		panic(err)
	}

	kid := uuid.UUID(hash.Sum(nil)[:16]).String()
	return kid
}
