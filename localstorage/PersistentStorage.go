package localstorage

import (
	"bytes"
	"clave/constants"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type PersistentStore struct {
	db *bolt.DB
}

var (
	ErrNotFound = errors.New("skv: key not found")

	ErrBadValue = errors.New("skv: bad value")

	bucketName = []byte("kv")

	persistentStorageSync sync.Once
	persistentInstance    *PersistentStore
)

func GetPersistentStorage() *PersistentStore {
	log.Println("Get Persistent Storage")
	persistentStorageSync.Do(func() {
		opts := &bolt.Options{
			Timeout: 50 * time.Millisecond,
		}
		if db, err := bolt.Open(constants.SecureVaultDB, 0640, opts); err != nil {

			if err != nil {
				log.Println("[GetPersistentStorage][Error] " + err.Error())
			}

			return
		} else {
			err := db.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists(bucketName)

				if err != nil {
					log.Println("[GetPersistentStorage][Error] " + err.Error())
				}

				return err
			})
			if err != nil {
				log.Println("[GetPersistentStorage][Error] " + err.Error())
				return
			} else {
				persistentInstance = &PersistentStore{db: db}
			}
		}

	})

	return persistentInstance

}

func (kvs *PersistentStore) SetValue(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		log.Println("[Put][Encode][Error] ", err)
		return err
	}

	text := buf.Bytes()
	password := []byte(constants.CommonDatabasePassword)

	// generate a new aes cipher using our 32 byte long password
	c, err := aes.NewCipher(password)
	// if there are any errors, handle them
	if err != nil {
		log.Println("[Put][Cipher][Error] " + err.Error())
		return err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric password cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		log.Println("[Put][GCM][Error] " + err.Error())
		return err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if numberOfBytesRead, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("[Put][Encrypt][Error] " + err.Error())
		return err
	} else {

		log.Println("[Encrypt] Read this many bytes --> ", numberOfBytesRead)
	}

	encryptedData := gcm.Seal(nonce, nonce, text, nil)

	return kvs.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), encryptedData)
	})
}

func (kvs *PersistentStore) Get(key string, value interface{}) error {
	log.Println("Get : ", key)
	return kvs.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, v := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		} else if value == nil {
			return nil
		} else {

			// Get the value, decrypt and send back

			key := []byte(constants.CommonDatabasePassword)
			ciphertext := v

			c, err := aes.NewCipher(key)
			if err != nil {
				log.Println("[Decrypt][Error] ", err)
				return err
			}

			gcm, err := cipher.NewGCM(c)
			if err != nil {
				log.Println("[Decrypt][Error] ", err)
				return err
			}

			nonceSize := gcm.NonceSize()
			if len(ciphertext) < nonceSize {
				log.Println("[Decrypt][Error] ", err)
				// return err
			}

			nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
			decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				log.Println("[Decrypt][Error] ", err)
				return err
			}

			d := gob.NewDecoder(bytes.NewReader(decryptedData))

			return d.Decode(value)

		}
	})
}

// Delete the entry with the given key. If no such key is present in the store,
// it returns ErrNotFound.
//
//	store.Delete("key42")
func (kvs *PersistentStore) DeleteKey(key string) error {
	log.Println("Delete Key : ", key)
	return kvs.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		} else {
			return c.Delete()
		}
	})
}

// Close closes the key-value store file.
func (kvs *PersistentStore) Close() error {
	log.Println("Closing PersistentStore")
	return kvs.db.Close()
}
