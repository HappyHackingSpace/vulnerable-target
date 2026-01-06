// Package store provides generic storage interfaces and factory functions for different storage backends
package store

import (
	"errors"

	"github.com/happyhackingspace/vt/pkg/store/disk"
)

// Type represents the type of storage backend
type Type string

const (
	// DiskStoreType indicates file-based disk storage
	DiskStoreType Type = "disk"
)

// Storage provides generic storage operations for storable objects
type Storage[T any] interface {
	Set(string, T) error
	Get(string) (T, error)
	GetAll() ([]T, error)
	Delete(string) error
	Close() error
}

// NewStorage creates a new storage instance based on the specified store type and configuration
func NewStorage[T any](storeType Type, config any) (Storage[T], error) {
	switch storeType {
	case DiskStoreType:
		cfg, ok := config.(*disk.Config)
		if !ok {
			return nil, errors.New("invalid config type for disk store")
		}
		return disk.NewStorageStore[T](cfg)
	default:
		return nil, errors.New("invalid store type")
	}
}
