package token

import "time"

// Maker is an interface for managing payloads
type Maker interface {
	// Create token creates token for the user
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken verifies if the token is valid
	VerifyToken(token string) (*Payload, error)
}
