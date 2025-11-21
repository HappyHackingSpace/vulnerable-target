package disk

import (
	"os"
	"path/filepath"
	"testing"

	bolt "go.etcd.io/bbolt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUser is a sample struct for testing
type TestUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// setupTestStore creates a temporary store for testing
func setupTestStore(t *testing.T) (*Store[TestUser], string) {
	t.Helper()

	// Create temporary directory for test database
	tmpDir := t.TempDir()

	config := &Config{
		FileName:   "test.db",
		BucketName: "test_bucket",
	}

	store, err := newStorageStoreWithPath[TestUser](config, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, store)

	t.Cleanup(func() {
		err = store.Close()
		require.Nil(t, err)
	})

	return store, tmpDir
}

// newStorageStoreWithPath creates a new disk-based storage instance with a custom base path (for testing)
func newStorageStoreWithPath[T any](config *Config, basePath string) (*Store[T], error) {
	vtHomeDir := filepath.Join(basePath, ".vt")
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

func TestNewStorageStore(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		FileName:   "test.db",
		BucketName: "users",
	}

	store, err := newStorageStoreWithPath[TestUser](config, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, store)
	defer func() {
		err = store.Close()
		require.Nil(t, err)
	}()

	// Verify the database file was created
	dbPath := filepath.Join(tmpDir, ".vt", "test.db")
	_, err = os.Stat(dbPath)
	assert.NoError(t, err, "database file should exist")

	// Verify config is set
	assert.Equal(t, config, store.Config)
}

func TestStore_Set(t *testing.T) {
	store, _ := setupTestStore(t)

	user := TestUser{
		ID:   "1",
		Name: "John Doe",
	}

	err := store.Set("user1", user)
	assert.NoError(t, err)
}

func TestStore_Get(t *testing.T) {
	store, _ := setupTestStore(t)

	expectedUser := TestUser{
		ID:   "1",
		Name: "John Doe",
	}

	// First, set the user
	err := store.Set("user1", expectedUser)
	require.NoError(t, err)

	// Now retrieve it
	actualUser, err := store.Get("user1")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)
}

func TestStore_GetAll(t *testing.T) {
	store, _ := setupTestStore(t)

	users := []TestUser{
		{ID: "1", Name: "John Doe"},
		{ID: "2", Name: "Jane Smith"},
		{ID: "3", Name: "Bob Johnson"},
	}

	// Store multiple users
	for i, user := range users {
		err := store.Set(user.ID, user)
		require.NoError(t, err, "failed to set user %d", i)
	}

	// Retrieve all users
	results, err := store.GetAll()
	assert.NoError(t, err)
	assert.Len(t, results, len(users))

	// Verify all users are present (order might differ)
	resultMap := make(map[string]TestUser)
	for _, user := range results {
		resultMap[user.ID] = user
	}

	for _, expectedUser := range users {
		actualUser, exists := resultMap[expectedUser.ID]
		assert.True(t, exists, "user %s should exist", expectedUser.ID)
		assert.Equal(t, expectedUser, actualUser)
	}
}

func TestStore_GetAll_EmptyBucket(t *testing.T) {
	store, _ := setupTestStore(t)

	// GetAll on empty bucket should return empty slice
	results, err := store.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestStore_Delete(t *testing.T) {
	store, _ := setupTestStore(t)

	user := TestUser{
		ID:   "1",
		Name: "John Doe",
	}

	// Set a user
	err := store.Set("user1", user)
	require.NoError(t, err)

	// Verify it exists
	_, err = store.Get("user1")
	require.NoError(t, err)

	// Delete the user
	err = store.Delete("user1")
	assert.NoError(t, err)

	// Verify it's gone
	_, err = store.Get("user1")
	assert.Error(t, err, "should return error for deleted key")
}

func TestStore_SetAndGetMultiple(t *testing.T) {
	store, _ := setupTestStore(t)

	users := map[string]TestUser{
		"user1": {ID: "1", Name: "Alice"},
		"user2": {ID: "2", Name: "Bob"},
		"user3": {ID: "3", Name: "Charlie"},
	}

	// Set multiple users
	for key, user := range users {
		err := store.Set(key, user)
		require.NoError(t, err)
	}

	// Get each user and verify
	for key, expectedUser := range users {
		actualUser, err := store.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, actualUser)
	}
}

func TestStore_UpdateExistingValue(t *testing.T) {
	store, _ := setupTestStore(t)

	originalUser := TestUser{
		ID:   "1",
		Name: "John Doe",
	}

	updatedUser := TestUser{
		ID:   "1",
		Name: "John Doe Jr.",
	}

	// Set original value
	err := store.Set("user1", originalUser)
	require.NoError(t, err)

	// Update the value
	err = store.Set("user1", updatedUser)
	assert.NoError(t, err)

	// Verify the update
	actualUser, err := store.Get("user1")
	assert.NoError(t, err)
	assert.Equal(t, updatedUser, actualUser)
	assert.NotEqual(t, originalUser, actualUser)
}

func TestStore_Close(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		FileName:   "test.db",
		BucketName: "test",
	}

	store, err := newStorageStoreWithPath[TestUser](config, tmpDir)
	require.NoError(t, err)

	// Close should not return an error
	err = store.Close()
	assert.NoError(t, err)
}

func TestStore_WithDifferentTypes(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("string values", func(t *testing.T) {
		config := &Config{
			FileName:   "string_test.db",
			BucketName: "strings",
		}

		store, err := newStorageStoreWithPath[string](config, tmpDir)
		require.NoError(t, err)
		defer func() {
			err = store.Close()
			require.Nil(t, err)
		}()

		err = store.Set("key1", "hello world")
		require.NoError(t, err)

		value, err := store.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, "hello world", value)
	})

	t.Run("int values", func(t *testing.T) {
		config := &Config{
			FileName:   "int_test.db",
			BucketName: "integers",
		}

		store, err := newStorageStoreWithPath[int](config, tmpDir)
		require.NoError(t, err)
		defer func() {
			err = store.Close()
			require.Nil(t, err)
		}()

		err = store.Set("key1", 42)
		require.NoError(t, err)

		value, err := store.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("map values", func(t *testing.T) {
		config := &Config{
			FileName:   "map_test.db",
			BucketName: "maps",
		}

		store, err := newStorageStoreWithPath[map[string]string](config, tmpDir)
		require.NoError(t, err)
		defer func() {
			err = store.Close()
			require.Nil(t, err)
		}()

		testMap := map[string]string{
			"foo": "bar",
			"baz": "qux",
		}

		err = store.Set("key1", testMap)
		require.NoError(t, err)

		value, err := store.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, testMap, value)
	})
}

func TestStore_ConcurrentOperations(t *testing.T) {
	store, _ := setupTestStore(t)

	// BoltDB handles concurrent reads well
	users := []TestUser{
		{ID: "1", Name: "User1"},
		{ID: "2", Name: "User2"},
		{ID: "3", Name: "User3"},
	}

	// Set users sequentially
	for _, user := range users {
		err := store.Set(user.ID, user)
		require.NoError(t, err)
	}

	// Get all users
	results, err := store.GetAll()
	assert.NoError(t, err)
	assert.Len(t, results, len(users))
}
