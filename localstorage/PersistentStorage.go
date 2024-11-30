package localstorage

import (
	"bytes"
	"clave/constants"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
)

var (
	ErrNotFound = errors.New("storage: key not found")
	ErrBadValue = errors.New("storage: bad value")

	storageSync sync.Once
	instance    *PersistentStore
)

type PersistentStore struct {
	db *badger.DB
}

func GetPersistentStorage() *PersistentStore {
	storageSync.Do(func() {

		dbPath := constants.SecureVaultDB
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Printf("[Storage] Failed to create database directory: %v", err)
			panic(err)
		}
		key := sha256.Sum256([]byte(constants.CommonDatabasePassword))
		opts := badger.DefaultOptions(dbPath).
			WithLogger(nil).
			WithCompression(options.ZSTD).
			WithEncryptionKey(key[:]).
			WithIndexCacheSize(100 << 20)

		db, err := badger.Open(opts)
		if err != nil {
			log.Printf("[Storage] Failed to open DB: %v", err)
			panic(err)
		}

		instance = &PersistentStore{db: db}
		go instance.runGC()
	})

	if instance == nil {
		panic("Failed to initialize storage")
	}

	return instance
}

func (ps *PersistentStore) SetValue(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}

	encrypted, err := ps.encrypt(buf.Bytes())
	if err != nil {
		return err
	}

	return ps.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), encrypted)
	})
}

func (ps *PersistentStore) Get(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}

	var data []byte
	err := ps.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			data = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		return err
	}

	decrypted, err := ps.decrypt(data)
	if err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewReader(decrypted)).Decode(value)
}

func (ps *PersistentStore) DeleteKey(key string) error {
	return ps.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (ps *PersistentStore) Close() error {
	if ps.db != nil {
		return ps.db.Close()
	}
	return nil
}

func (ps *PersistentStore) encrypt(data []byte) ([]byte, error) {
	key := sha256.Sum256([]byte(constants.CommonDatabasePassword))

	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (ps *PersistentStore) decrypt(data []byte) ([]byte, error) {
	key := sha256.Sum256([]byte(constants.CommonDatabasePassword))

	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("malformed data")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (ps *PersistentStore) IsHealthy() bool {
	if ps.db == nil {
		return false
	}

	err := ps.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("health_check"))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		return err
	})

	return err == nil
}

func (ps *PersistentStore) runGC() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Storage] GC panic recovered: %v", r)
				}
			}()

			err := ps.db.RunValueLogGC(0.5)
			if err != nil && err != badger.ErrNoRewrite {
				log.Printf("[Storage] GC error: %v", err)
			}
		}()
	}
}

func (ps *PersistentStore) HasKey(key string) bool {
	err := ps.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	return err == nil
}

func (ps *PersistentStore) Encrypt(data []byte) ([]byte, error) {
	return ps.encrypt(data)
}

func (ps *PersistentStore) Decrypt(data []byte) ([]byte, error) {
	return ps.decrypt(data)
}
