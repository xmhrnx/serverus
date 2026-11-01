package domain

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID returns a random identifier for a new server, snippet or tunnel.
func NewID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
