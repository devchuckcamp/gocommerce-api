package utils

import "github.com/google/uuid"

// GenerateID generates a new UUID v4 ID
func GenerateID() string {
	return uuid.New().String()
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	return "ORD-" + uuid.New().String()[:8]
}
