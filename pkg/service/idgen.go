package service

import (
	"strings"

	"github.com/google/uuid"
)

type IDGenerator interface {
	NewID() (string, error)
}

type uuidGen struct {
	length int
}

func NewUUIDGen(length int) *uuidGen {
	return &uuidGen{length: length}
}

func (g *uuidGen) NewID() (string, error) {
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	if len(id) > g.length {
		return id[:g.length], nil
	}
	return id, nil
}
