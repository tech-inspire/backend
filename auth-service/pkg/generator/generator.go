package generator

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/matoous/go-nanoid/v2"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

func (Generator) GenerateString(length int) string {
	return gonanoid.Must(length)
}

func (Generator) NewUUID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func NumberCode(length int) string {
	buf := make([]byte, length)

	for i := range buf {
		buf[i] = '0' + byte(rand.Intn(9))
	}

	return string(buf)
}
