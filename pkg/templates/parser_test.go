package templates

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTemplate(t *testing.T) {
	// create temp dir
	tempDir := t.TempDir()

	// create dummy content
	templateContent := `
id: example-template

info:
  name: Vulnerable Target
  author: hhsteam
  description: |
    Vulnerable Target
  references:
    - http://www.vulnerabletarget.com
  technologies:
    - php
    - mysql
  tags:
    - owasp
    - web
    - vulnerabilities
  metadata:

providers:
  online:
    targets:
      - vulnerabletarget.com

  aws:
    path:
`
	err := os.WriteFile(filepath.Join(tempDir, "index.yaml"), []byte(templateContent), 0644)
	assert.NoError(t, err)

	tpl, err := LoadTemplate(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, "example-template", tpl.ID)
	assert.Equal(t, "Vulnerable Target", tpl.Info.Name)
	assert.Equal(t, "hhsteam", tpl.Info.Author)
	assert.Equal(t, 1, len(tpl.Info.References))
	assert.Equal(t, 3, len(tpl.Info.Tags))
	assert.Equal(t, 3, len(tpl.Info.Tags))
	assert.Contains(t, tpl.Providers, "aws")

	// case of none exist path
	_, err = LoadTemplate("/non/existent/path")
	assert.Error(t, err)

}
