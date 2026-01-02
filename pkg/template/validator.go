package template

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	templateIDRegex     = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
	allowedProviderExts = map[string]bool{
		".yml":  true,
		".yaml": true,
	}
)

// Validate validates the template structure and content.
func (template Template) Validate() error {
	if template.ID == "" {
		return fmt.Errorf("id can not be empty")
	}

	if !templateIDRegex.MatchString(template.ID) {
		return fmt.Errorf("template '%s': id contains invalid characters", template.ID)
	}

	if len(template.Providers) == 0 {
		return fmt.Errorf("template '%s': no providers specified in the template", template.ID)
	}

	for k, v := range template.Providers {
		pcError := v.Validate(template.ID, k)
		if pcError != nil {
			return pcError
		}
	}

	infoError := template.Info.Validate(template.ID)
	if infoError != nil {
		return infoError
	}

	return nil
}

// Validate validates the info structure and content.
func (info Info) Validate(templateID string) error {
	if info.Name == "" {
		return fmt.Errorf("template '%s': name can not be empty", templateID)
	}
	if info.Author == "" {
		return fmt.Errorf("template '%s': author can not be empty", templateID)
	}
	if len(info.Targets) == 0 {
		return fmt.Errorf("template '%s': targets can not be empty", templateID)
	}
	if info.Type == "" {
		return fmt.Errorf("template '%s': type can not be empty", templateID)
	}
	if len(info.Tags) == 0 {
		return fmt.Errorf("template '%s': tags can not be empty", templateID)
	}
	return nil
}

// Validate validates the provider configuration structure and content.
func (pc ProviderConfig) Validate(templateID, name string) error {
	providerPath := pc.Path
	if providerPath == "" {
		return fmt.Errorf("template '%s', provider '%s': path is empty", templateID, name)
	}

	if filepath.IsAbs(providerPath) {
		return fmt.Errorf("template '%s', provider '%s': absolute paths are not allowed", templateID, name)
	}

	if strings.Contains(providerPath, "..") {
		return fmt.Errorf("template '%s', provider '%s': path contains invalid '..' segments", templateID, name)
	}

	ext := filepath.Ext(providerPath)
	if !isAllowedExtension(ext) {
		return fmt.Errorf("template '%s', provider '%s': provider file must have one of the allowed extensions: %v", templateID, name, allowedProviderExts)
	}

	return nil
}

func isAllowedExtension(ext string) bool {
	return allowedProviderExts[ext]
}
