package token

import (
	"time"
)

// Maker is an interface for making tokens
type Maker interface {
	// CreateToken creates a new token for a specific user and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the input token is valid
	VerifyToken(token string) (*Payload, error)
}
