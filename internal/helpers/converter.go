package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

func StringToUUID(id string) (uuid.UUID, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return uuidID, nil
}
