package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateTempFile(content string, name string) (string, error) {
	tempDir := filepath.Join(os.TempDir(), "vt-folder")

	err := os.MkdirAll(tempDir, 0o700)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	filePath := filepath.Join(tempDir, name)

	err = os.WriteFile(filePath, []byte(content), 0o600)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func DeleteFile(path string) error {
	return os.RemoveAll(path)
}
