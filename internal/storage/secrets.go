package storage

import (
	"errors"
	"fmt"
)

// SecretStore keeps credentials in the platform's keychain. The app
// passes giode's SecureStore; tests pass a map.
type SecretStore interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

// Secrets names the credentials the app stores for each server.
type Secrets struct {
	Store SecretStore
}

// ErrSecretNotFound is returned by a SecretStore that holds nothing
// under the asked key.
var ErrSecretNotFound = errors.New("secret not found")

// PasswordKey is the key of a server's password.
func PasswordKey(serverID string) string { return fmt.Sprintf("server/%s/password", serverID) }

// PassphraseKey is the key of a server's private key passphrase.
func PassphraseKey(serverID string) string { return fmt.Sprintf("server/%s/passphrase", serverID) }

// Get returns the secret, or "" when none is stored.
func (s Secrets) Get(key string) string {
	if s.Store == nil {
		return ""
	}
	v, err := s.Store.Get(key)
	if err != nil {
		return ""
	}
	return string(v)
}

// Set stores the secret; an empty value deletes it.
func (s Secrets) Set(key, value string) error {
	if s.Store == nil {
		return nil
	}
	if value == "" {
		return s.Delete(key)
	}
	return s.Store.Set(key, []byte(value))
}
