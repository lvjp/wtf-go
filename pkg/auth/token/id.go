package token

import (
	"fmt"

	"github.com/google/uuid"
)

type IDGenerator func() (string, error)

func UUIDGenerator() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("UUID generation error: %v", err)
	}

	return token.String(), nil
}
