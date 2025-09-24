// Package state provides deployment management
package state

import (
	"fmt"
	"time"

	"github.com/happyhackingspace/vulnerable-target/pkg/store"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/disk"
	"github.com/happyhackingspace/vulnerable-target/pkg/store/storable"
)

// Deployment represents the status of an environment on a specified provider
type Deployment struct {
	ProviderName string
	TemplateID   string
	Status       string
	storable.Struct
}

// Manager provides storage operations for deployments
type Manager struct {
	store store.Storage[Deployment]
}

// NewManager creates a new manager with pre-defined disk storage configuration
func NewManager() (*Manager, error) {
	store, err := store.NewStorage[Deployment](store.DiskStoreType, disk.Config{
		FileName:   "deployments.db",
		BucketName: "deployment",
	})
	if err != nil {
		return nil, err
	}
	return &Manager{store: store}, nil
}

// AddNewDeployment creates a new deployment record with running status
func (m Manager) AddNewDeployment(providerName, templateID string) error {
	deployment := Deployment{ProviderName: providerName, TemplateID: templateID, Status: "running", Struct: storable.Struct{CreatedAt: time.Now()}}
	err := m.store.Set(fmt.Sprintf("%s:%s", deployment.ProviderName, deployment.TemplateID), deployment)
	return err
}

// RemoveDeployment deletes a deployment record by provider name and template ID
func (m Manager) RemoveDeployment(providerName, templateID string) error {
	err := m.store.Delete(fmt.Sprintf("%s:%s", providerName, templateID))
	return err
}

// DeploymentExist checks if a deployment exists for the given provider and template
func (m Manager) DeploymentExist(providerName, templateID string) (bool, error) {
	_, err := m.store.Get(fmt.Sprintf("%s:%s", providerName, templateID))
	if err != nil {
		return false, err
	}
	return true, err
}

// ListDeployments returns all deployment records from storage
func (m Manager) ListDeployments() ([]Deployment, error) {
	deployments, err := m.store.GetAll()
	return deployments, err
}
