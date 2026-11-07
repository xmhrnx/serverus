package storage

import (
	"errors"
)

// Delete removes the secret. Deleting a missing secret is not an error.
func (s Secrets) Delete(key string) error {
	if s.Store == nil {
		return nil
	}
	err := s.Store.Delete(key)
	if err != nil && errors.Is(err, ErrSecretNotFound) {
		return nil
	}
	return err
}

// ForgetServer deletes every secret of a server.
func (s Secrets) ForgetServer(serverID string) {
	_ = s.Delete(PasswordKey(serverID))
	_ = s.Delete(PassphraseKey(serverID))
}
