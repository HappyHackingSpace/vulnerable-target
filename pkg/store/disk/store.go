// Package disk provides configuration for disk-based storage implementations
package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// Store provides disk-based storage using BoltDB for storable objects
type Store[T any] struct {
	Config *Config
	db     *bolt.DB
}

// Set stores a key-value pair in the database
func (s *Store[T]) Set(key string, value T) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.Config.BucketName))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		return bucket.Put([]byte(key), data)
	})
}

// Get retrieves a value by key from the database
func (s *Store[T]) Get(key string) (T, error) {
	var zero T
	var result T

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.Config.BucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", s.Config.BucketName)
		}
		data := bucket.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key %s not found", key)
		}

		return json.Unmarshal(data, &result)
	})

	if err != nil {
		return zero, err
	}
	return result, nil
}

// GetAll retrieves all values from the database bucket
func (s *Store[T]) GetAll() ([]T, error) {
	var results []T

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.Config.BucketName))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var item T
			if err := json.Unmarshal(v, &item); err != nil {
				return fmt.Errorf("failed to unmarshal item with key %s: %w", string(k), err)
			}
			results = append(results, item)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return results, nil
}

// Delete removes a key-value pair from the database
func (s *Store[T]) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.Config.BucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", s.Config.BucketName)
		}

		return bucket.Delete([]byte(key))
	})
}

// Close closes the database connection
func (s *Store[T]) Close() error {
	return s.db.Close()
}

// NewStorageStore creates a new disk-based storage instance with the given configuration
func NewStorageStore[T any](config *Config) (*Store[T], error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	vtHomeDir := filepath.Join(userHomeDir, ".vt")
	if err := os.MkdirAll(vtHomeDir, 0750); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(vtHomeDir, config.FileName)
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Store[T]{
		db:     db,
		Config: config,
	}, nil
}
