package generator

import (
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
