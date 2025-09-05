package templates

import (
	"os"
	"path"

	yaml "gopkg.in/yaml.v3"
)

func LoadTemplate(filepath string) (Template, error) {
	var template Template
	file, err := os.ReadFile(path.Join(filepath, "index.yaml")) // #nosec G304
	if err != nil {
		return template, err
	}
	err = yaml.Unmarshal(file, &template)
	if err != nil {
		return template, err
	}
	err = template.Validate()
	return template, err
}
