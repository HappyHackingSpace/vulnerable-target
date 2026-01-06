package template

import (
	"fmt"
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
