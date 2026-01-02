package template

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	Templates = map[string]Template{
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
	firstTemplate, err := GetByID("example-template-1")
	assert.Nil(t, err)
	assert.Equal(t, "example-template-1", firstTemplate.ID)
	noneExistTemplateID := "none-exist-template"
	noneExistingTemplate, err := GetByID(noneExistTemplateID)
	assert.Nil(t, noneExistingTemplate)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("template %s not found", noneExistTemplateID))

}
