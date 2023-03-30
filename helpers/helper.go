package helpers

import "github.com/google/uuid"

func GenerateRandomUUID() string {
	// Generate a new UUID
	id := uuid.New()

	return id.String()
}
