package bootstrap

import (
	"errors"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/storage"
)

// secureStore adapts giode's keychain access to storage.SecretStore.
type secureStore struct{}

func (secureStore) Get(key string) ([]byte, error) {
	v, err := giode.SecureStore.Get(key)
	if errors.Is(err, giode.ErrNotFound) {
		return nil, storage.ErrSecretNotFound
	}
	return v, err
}

func (secureStore) Set(key string, value []byte) error { return giode.SecureStore.Set(key, value) }

func (secureStore) Delete(key string) error {
	err := giode.SecureStore.Delete(key)
	if errors.Is(err, giode.ErrNotFound) {
		return storage.ErrSecretNotFound
	}
	return err
}
