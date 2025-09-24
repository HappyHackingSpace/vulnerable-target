// Package storable provides interfaces and base implementations for objects that can be stored with timestamps
package storable

import "time"

// Interface defines the contract for objects that can be stored with creation timestamps
type Interface interface {
	GetCreatedAt() time.Time
}

// Struct provides a base implementation of the Interface with a creation timestamp field
type Struct struct {
	CreatedAt time.Time
}

// GetCreatedAt returns the creation timestamp of the storable object
func (s Struct) GetCreatedAt() time.Time {
	return s.CreatedAt
}
