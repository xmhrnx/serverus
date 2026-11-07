package storage

import (
	"testing"
)

type mapStore map[string][]byte

func (m mapStore) Get(k string) ([]byte, error) {
	v, ok := m[k]
	if !ok {
		return nil, ErrSecretNotFound
	}
	return v, nil
}

func (m mapStore) Set(k string, v []byte) error { m[k] = v; return nil }

func (m mapStore) Delete(k string) error {
	if _, ok := m[k]; !ok {
		return ErrSecretNotFound
	}
	delete(m, k)
	return nil
}

func TestSecrets(t *testing.T) {
	s := Secrets{Store: mapStore{}}
	s.Set(PasswordKey("a"), "pw")
	if s.Get(PasswordKey("a")) != "pw" {
		t.Fatal("stored password not returned")
	}
	s.ForgetServer("a")
	if s.Get(PasswordKey("a")) != "" {
		t.Fatal("forgotten password still returned")
	}
	if err := s.Delete("missing"); err != nil {
		t.Fatalf("deleting a missing secret: %v", err)
	}
}
