package templates

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func LoadTemplate(filepath string) (Template, error) {
	var tmpl Template
	file, err := os.ReadFile(filepath)
	if err != nil {
		return tmpl, err
	}
	err = yaml.Unmarshal(file, &tmpl)
	return tmpl, err
}

func ValidateTemplate(tmpl Template) error {
	if tmpl.ID == "" {
		return fmt.Errorf("template ID is missing")
	}
	if len(tmpl.Providers) == 0 {
		return fmt.Errorf("no providers specified in the template")
	}
	return nil
}
