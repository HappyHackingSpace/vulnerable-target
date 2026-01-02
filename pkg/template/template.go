package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
)

// TemplateRemoteRepository is a constant for repo url.
const TemplateRemoteRepository string = "https://github.com/HappyHackingSpace/vt-templates"

// Template represents a vulnerable target environment configuration.
type Template struct {
	ID             string                    `yaml:"id"`
	Info           Info                      `yaml:"info"`
	ProofOfConcept map[string][]string       `yaml:"poc"`
	Remediation    []string                  `yaml:"remediation"`
	Providers      map[string]ProviderConfig `yaml:"providers"`
	PostInstall    []string                  `yaml:"post-install"`
}

// Info contains metadata about a template.
type Info struct {
	Name             string   `yaml:"name"`
	Description      string   `yaml:"description"`
	Author           string   `yaml:"author"`
	Targets          []string `yaml:"targets"`
	Type             string   `yaml:"type"`
	AffectedVersions []string `yaml:"affected_versions"`
	FixedVersion     string   `yaml:"fixed_version"`
	Cwe              string   `yaml:"cwe"`
	Cvss             Cvss     `yaml:"cvss"`
	Tags             []string `yaml:"tags"`
	References       []string `yaml:"references"`
}

// ProviderConfig contains configuration for a specific provider.
type ProviderConfig struct {
	Path string `yaml:"path"`
}

// Cvss represents Common Vulnerability Scoring System information.
type Cvss struct {
	Score   string `yaml:"score"`
	Metrics string `yaml:"metrics"`
}

// LoadTemplates loads all templates from the given repository path.
// If the repository doesn't exist, it clones it first.
// Returns a map of templates indexed by their ID.
func LoadTemplates(repoPath string) (map[string]Template, error) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Info().Msg("Fetching templates for the first time")
		if err := cloneTemplatesRepo(repoPath, false); err != nil {
			return nil, fmt.Errorf("failed to clone templates repository: %w", err)
		}
	}

	return loadTemplatesFromDirectory(repoPath)
}

// loadTemplatesFromDirectory reads all templates from the given path.
// Returns a map of templates indexed by their ID.
func loadTemplatesFromDirectory(repoPath string) (map[string]Template, error) {
	templates := make(map[string]Template)

	dirEntry, err := os.ReadDir(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", repoPath, err)
	}

	for _, categoryEntry := range dirEntry {
		if strings.HasPrefix(categoryEntry.Name(), ".") || !categoryEntry.IsDir() {
			continue
		}

		categoryPath := filepath.Join(repoPath, categoryEntry.Name())
		categoryTemplates, err := loadTemplatesFromCategory(categoryPath, categoryEntry.Name())
		if err != nil {
			return nil, err
		}

		for id, tmpl := range categoryTemplates {
			templates[id] = tmpl
		}
	}

	return templates, nil
}

// loadTemplatesFromCategory loads all templates within a single category directory.
// Returns a map of templates indexed by their ID.
func loadTemplatesFromCategory(categoryPath, categoryName string) (map[string]Template, error) {
	templates := make(map[string]Template)

	templateEntries, err := os.ReadDir(categoryPath)
	if err != nil {
		return nil, fmt.Errorf("error reading category %s: %w", categoryName, err)
	}

	for _, entry := range templateEntries {
		if strings.HasPrefix(entry.Name(), ".") || !entry.IsDir() {
			continue
		}

		templatePath := filepath.Join(categoryPath, entry.Name())
		tmpl, err := LoadTemplate(templatePath)
		if err != nil {
			return nil, fmt.Errorf("error loading template %s: %w", entry.Name(), err)
		}
		if tmpl.ID != entry.Name() {
			return nil, fmt.Errorf("template id '%s' and directory name '%s' should match", tmpl.ID, entry.Name())
		}
		templates[tmpl.ID] = tmpl
	}

	return templates, nil
}

// SyncTemplates downloads or updates all templates from the remote repository.
func SyncTemplates(repoPath string) error {
	log.Info().Msgf("cloning %s", TemplateRemoteRepository)
	if err := cloneTemplatesRepo(repoPath, true); err != nil {
		return fmt.Errorf("failed to sync templates: %w", err)
	}
	return nil
}

// ListTemplates displays all available templates in a table format.
func ListTemplates(templates map[string]Template) {
	ListTemplatesWithFilter(templates, "")
}

// ListTemplatesWithFilter displays templates in a table format, optionally filtered by tag.
func ListTemplatesWithFilter(templates map[string]Template, filterTag string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Author", "Targets", "Type", "Tags"})

	count := 0
	for _, tmpl := range templates {
		if filterTag != "" {
			hasTag := false
			for _, tag := range tmpl.Info.Tags {
				if strings.EqualFold(tag, filterTag) || strings.Contains(strings.ToLower(tag), strings.ToLower(filterTag)) {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		tags := strings.Join(tmpl.Info.Tags, ", ")
		targets := strings.Join(tmpl.Info.Targets, ", ")
		t.AppendRow(table.Row{
			tmpl.ID,
			tmpl.Info.Name,
			tmpl.Info.Author,
			targets,
			tmpl.Info.Type,
			tags,
		})
		count++
	}

	if count == 0 {
		if filterTag != "" {
			fmt.Printf("No templates found with tag matching '%s'\n", filterTag)
		} else {
			fmt.Println("No templates found")
		}
		return
	}

	if filterTag != "" {
		t.SetCaption("Found %d templates with tag matching '%s'", count, filterTag)
	} else {
		t.SetCaption("there are %d templates", count)
	}
	t.SetIndexColumn(0)
	t.Render()
}

// GetByID retrieves a template by its ID from the given templates map.
func GetByID(templates map[string]Template, templateID string) (*Template, error) {
	tmpl, ok := templates[templateID]
	if !ok || tmpl.ID == "" {
		return nil, fmt.Errorf("template %s not found", templateID)
	}
	return &tmpl, nil
}
