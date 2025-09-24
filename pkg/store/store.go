// Package store provides generic storage interfaces and factory functions for different storage backends
package store

import (
	"errors"

	"github.com/happyhackingspace/vulnerable-target/pkg/store/config"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/disk"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/storable"
)

// Type represents the type of storage backend
type Type string

const (
	// DiskStoreType indicates file-based disk storage
	DiskStoreType Type = "disk"
)

// Storage provides generic storage operations for storable objects
type Storage[T storable.Interface] interface {
	Set(string, T) error
	Get(string) (T, error)
	GetAll() ([]T, error)
	Delete(string) error
	Close() error
}

// NewStorage creates a new storage instance based on the specified store type and configuration
func NewStorage[T storable.Interface](storeType Type, config config.Interface) (Storage[T], error) {
	switch storeType {
	case DiskStoreType:
		if diskConfig, ok := config.(disk.Config); ok {
			return disk.NewStorageStore[T](diskConfig)
		}
		return nil, errors.New("invalid config type for disk store")
	default:
		return nil, errors.New("invalid store type")
	}
}
