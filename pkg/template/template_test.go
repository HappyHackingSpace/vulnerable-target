package template

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	templates := map[string]Template{
		"example-template-1": {
			ID: "example-template-1",
			Info: Info{
				Name:   "example-template-name",
				Author: "example-template-author",
				Tags: []string{
					"test",
					"example",
				},
			},
		},
	}
	firstTemplate, err := GetByID(templates, "example-template-1")
	assert.Nil(t, err)
	assert.Equal(t, "example-template-1", firstTemplate.ID)
	noneExistTemplateID := "none-exist-template"
	noneExistingTemplate, err := GetByID(templates, noneExistTemplateID)
	assert.Nil(t, noneExistingTemplate)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("template %s not found", noneExistTemplateID))
}

func TestGetByTags(t *testing.T) {
	templates := map[string]Template{
		"sqli-template": {
			ID: "sqli-template",
			Info: Info{
				Name: "SQL Injection Lab",
				Tags: []string{"sqli", "web", "owasp"},
			},
		},
		"xss-template": {
			ID: "xss-template",
			Info: Info{
				Name: "XSS Lab",
				Tags: []string{"xss", "web", "owasp"},
			},
		},
		"ssrf-template": {
			ID: "ssrf-template",
			Info: Info{
				Name: "SSRF Lab",
				Tags: []string{"ssrf", "web"},
			},
		},
	}

	// Test single tag match
	matched, err := GetByTags(templates, []string{"sqli"})
	assert.NoError(t, err)
	assert.Len(t, matched, 1)
	assert.Equal(t, "sqli-template", matched[0].ID)

	// Test multiple tags (OR logic)
	matched, err = GetByTags(templates, []string{"sqli", "xss"})
	assert.NoError(t, err)
	assert.Len(t, matched, 2)

	// Test tag that matches multiple templates
	matched, err = GetByTags(templates, []string{"web"})
	assert.NoError(t, err)
	assert.Len(t, matched, 3)

	// Test case-insensitive matching
	matched, err = GetByTags(templates, []string{"SQLI"})
	assert.NoError(t, err)
	assert.Len(t, matched, 1)

	// Test substring matching
	matched, err = GetByTags(templates, []string{"owa"})
	assert.NoError(t, err)
	assert.Len(t, matched, 2) // sqli-template and xss-template have "owasp"

	// Test no matches
	matched, err = GetByTags(templates, []string{"nonexistent"})
	assert.Error(t, err)
	assert.Nil(t, matched)
	assert.Contains(t, err.Error(), "no templates found matching tags")

	// Test empty tags
	matched, err = GetByTags(templates, []string{})
	assert.Error(t, err)
	assert.Nil(t, matched)
	assert.Contains(t, err.Error(), "no tags provided")

	// Test whitespace-only tags are ignored
	matched, err = GetByTags(templates, []string{"  ", "sqli"})
	assert.NoError(t, err)
	assert.Len(t, matched, 1)
}

func TestTemplateMatchesTags(t *testing.T) {
	tmpl := &Template{
		ID: "test-template",
		Info: Info{
			Tags: []string{"sqli", "XSS", "OWASP-Top10"},
		},
	}

	// Exact match
	assert.True(t, templateMatchesTags(tmpl, []string{"sqli"}))

	// Case-insensitive match
	assert.True(t, templateMatchesTags(tmpl, []string{"SQLI"}))
	assert.True(t, templateMatchesTags(tmpl, []string{"xss"}))

	// Substring match
	assert.True(t, templateMatchesTags(tmpl, []string{"owasp"}))
	assert.True(t, templateMatchesTags(tmpl, []string{"top10"}))

	// No match
	assert.False(t, templateMatchesTags(tmpl, []string{"nonexistent"}))

	// Empty filter tags
	assert.False(t, templateMatchesTags(tmpl, []string{}))
	assert.False(t, templateMatchesTags(tmpl, []string{"  ", ""}))
}

// createTestTemplate creates a template directory with an index.yaml file
func createTestTemplate(t *testing.T, basePath, templateID string) {
	t.Helper()
	templateContent := fmt.Sprintf(`
id: %s

info:
  name: Test Template %s
  author: testauthor
  description: Test description
  type: Lab
  targets:
    - test
  tags:
    - test

providers:
  docker-compose:
    path: "docker-compose.yaml"
`, templateID, templateID)

	templateDir := filepath.Join(basePath, templateID)
	err := os.MkdirAll(templateDir, 0750)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(templateDir, "index.yaml"), []byte(templateContent), 0644)
	assert.NoError(t, err)
}

func TestIsTemplateDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create a directory without index.yaml
	nonTemplateDir := filepath.Join(tempDir, "not-a-template")
	err := os.MkdirAll(nonTemplateDir, 0750)
	assert.NoError(t, err)
	assert.False(t, isTemplateDirectory(nonTemplateDir))

	// Create a directory with index.yaml
	templateDir := filepath.Join(tempDir, "is-a-template")
	err = os.MkdirAll(templateDir, 0750)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(templateDir, "index.yaml"), []byte("test"), 0644)
	assert.NoError(t, err)
	assert.True(t, isTemplateDirectory(templateDir))

	// Non-existent directory
	assert.False(t, isTemplateDirectory(filepath.Join(tempDir, "nonexistent")))
}

func TestLoadTemplatesFromCategoryNestedStructure(t *testing.T) {
	tempDir := t.TempDir()

	// Create a nested directory structure:
	// tempDir/
	//   ├── template-1/
	//   │   └── index.yaml
	//   └── subdir/
	//       ├── template-2/
	//       │   └── index.yaml
	//       └── nested/
	//           └── template-3/
	//               └── index.yaml

	createTestTemplate(t, tempDir, "template-1")
	subdir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subdir, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, subdir, "template-2")

	nested := filepath.Join(subdir, "nested")
	err = os.MkdirAll(nested, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, nested, "template-3")

	// Scan and verify
	templates, err := loadTemplatesFromCategory(tempDir, "test-category")
	assert.NoError(t, err)
	assert.Len(t, templates, 3)
	assert.Contains(t, templates, "template-1")
	assert.Contains(t, templates, "template-2")
	assert.Contains(t, templates, "template-3")
}

func TestLoadTemplatesFromCategorySkipsHiddenDirs(t *testing.T) {
	tempDir := t.TempDir()

	// Create a visible template
	createTestTemplate(t, tempDir, "visible-template")

	// Create a hidden directory with a template (should be skipped)
	hiddenDir := filepath.Join(tempDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, hiddenDir, "hidden-template")

	// Scan and verify only visible template is found
	templates, err := loadTemplatesFromCategory(tempDir, "test-category")
	assert.NoError(t, err)
	assert.Len(t, templates, 1)
	assert.Contains(t, templates, "visible-template")
	assert.NotContains(t, templates, "hidden-template")
}

func TestLoadTemplatesFromCategoryMaxDepth(t *testing.T) {
	tempDir := t.TempDir()

	// Create a deeply nested structure exceeding maxScanDepth
	currentPath := tempDir
	for i := 0; i <= maxScanDepth+1; i++ {
		currentPath = filepath.Join(currentPath, fmt.Sprintf("level%d", i))
		err := os.MkdirAll(currentPath, 0750)
		assert.NoError(t, err)
	}

	// Scan should fail due to max depth exceeded
	_, err := loadTemplatesFromCategory(tempDir, "test-category")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum directory depth")
}

func TestLoadTemplatesFromCategoryDuplicateID(t *testing.T) {
	tempDir := t.TempDir()

	// Create two templates with the same ID in different subdirectories
	createTestTemplate(t, tempDir, "duplicate-template")

	subdir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subdir, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, subdir, "duplicate-template")

	// Scan should fail due to duplicate template ID
	_, err = loadTemplatesFromCategory(tempDir, "test-category")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate template id")
}

func TestLoadTemplatesFromCategoryWithSubdirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create category with nested templates
	// category/
	//   ├── direct-template/
	//   │   └── index.yaml
	//   └── subcategory/
	//       └── nested-template/
	//           └── index.yaml

	createTestTemplate(t, tempDir, "direct-template")

	subcategory := filepath.Join(tempDir, "subcategory")
	err := os.MkdirAll(subcategory, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, subcategory, "nested-template")

	// Load and verify
	templates, err := loadTemplatesFromCategory(tempDir, "test-category")
	assert.NoError(t, err)
	assert.Len(t, templates, 2)
	assert.Contains(t, templates, "direct-template")
	assert.Contains(t, templates, "nested-template")
}

func TestLoadTemplatesFromDirectoryWithNestedCategories(t *testing.T) {
	tempDir := t.TempDir()

	// Create structure:
	// tempDir/
	//   ├── category1/
	//   │   └── template-a/
	//   │       └── index.yaml
	//   └── category2/
	//       ├── template-b/
	//       │   └── index.yaml
	//       └── subgroup/
	//           └── template-c/
	//               └── index.yaml

	category1 := filepath.Join(tempDir, "category1")
	err := os.MkdirAll(category1, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, category1, "template-a")

	category2 := filepath.Join(tempDir, "category2")
	err = os.MkdirAll(category2, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, category2, "template-b")

	subgroup := filepath.Join(category2, "subgroup")
	err = os.MkdirAll(subgroup, 0750)
	assert.NoError(t, err)
	createTestTemplate(t, subgroup, "template-c")

	// Load and verify
	templates, err := loadTemplatesFromDirectory(tempDir)
	assert.NoError(t, err)
	assert.Len(t, templates, 3)
	assert.Contains(t, templates, "template-a")
	assert.Contains(t, templates, "template-b")
	assert.Contains(t, templates, "template-c")
}
