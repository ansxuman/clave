package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	PIN_KEY  = "user_pin"
	SALT_KEY = "pin_salt"
)

type Service struct {
	storage Storage
}

type Storage interface {
	IsHealthy() bool
	HasKey(key string) bool
	SetValue(key string, value interface{}) error
	Get(key string, value interface{}) error
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) HasPin() bool {
	return s.storage.HasKey(PIN_KEY) && s.storage.HasKey(SALT_KEY)
}

func (s *Service) SetupPin(pin string) error {
	if len(pin) < 6 {
		return errors.New("PIN must be at least 6 characters")
	}
	if s.HasPin() {
		return errors.New("PIN is already set")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(pin), salt, 1, 64*1024, 4, 32)

	if err := s.storage.SetValue(SALT_KEY, base64.StdEncoding.EncodeToString(salt)); err != nil {
		return err
	}

	return s.storage.SetValue(PIN_KEY, base64.StdEncoding.EncodeToString(hash))
}

func (s *Service) VerifyPin(pin string) (bool, error) {
	if !s.storage.IsHealthy() {
		return false, errors.New("storage is not healthy")
	}

	var storedHashStr string
	var saltStr string

	if err := s.storage.Get(PIN_KEY, &storedHashStr); err != nil {
		return false, err
	}

	if err := s.storage.Get(SALT_KEY, &saltStr); err != nil {
		return false, err
	}

	storedHash, err := base64.StdEncoding.DecodeString(storedHashStr)
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(pin), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, storedHash) == 1, nil
}
