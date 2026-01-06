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

// String returns template fields as a table
func (t Template) String() string {
	tw := table.NewWriter()
	tw.AppendRow(table.Row{"ID", t.ID})
	tw.AppendRow(table.Row{"Name", t.Info.Name})
	tw.AppendRow(table.Row{"Description", t.Info.Description})
	tw.AppendRow(table.Row{"Author", t.Info.Author})
	tw.AppendRow(table.Row{"Type", t.Info.Type})
	tw.AppendRow(table.Row{"Targets", formatList(t.Info.Targets)})
	tw.AppendRow(table.Row{"Affected Versions", formatList(t.Info.AffectedVersions)})
	tw.AppendRow(table.Row{"Fixed Version", t.Info.FixedVersion})
	tw.AppendRow(table.Row{"CWE", t.Info.Cwe})
	tw.AppendRow(table.Row{"CVSS Score", t.Info.Cvss.Score})
	tw.AppendRow(table.Row{"CVSS Metrics", t.Info.Cvss.Metrics})
	tw.AppendRow(table.Row{"Tags", formatList(t.Info.Tags)})
	tw.AppendRow(table.Row{"References", formatList(t.Info.References)})
	tw.AppendRow(table.Row{"Proof of Concept", formatPoc(t.ProofOfConcept)})
	tw.AppendRow(table.Row{"Remediation", formatList(t.Remediation)})
	tw.AppendRow(table.Row{"Providers", formatProviders(t.Providers)})
	tw.AppendRow(table.Row{"Post Install", formatList(t.PostInstall)})

	tw.Style().Options.DrawBorder = true
	tw.Style().Options.SeparateRows = true
	tw.Style().Options.SeparateColumns = true

	return tw.Render()
}

func formatPoc(poc map[string][]string) string {
	if len(poc) == 0 {
		return ""
	}
	var parts []string
	for key, values := range poc {
		parts = append(parts, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
	}
	return strings.Join(parts, "\n")
}

func formatProviders(providers map[string]ProviderConfig) string {
	if len(providers) == 0 {
		return ""
	}
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return strings.Join(names, "\n")
}

func formatList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, "\n")
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

// maxScanDepth limits the depth of recursive directory scanning to prevent
// infinite loops from circular symlinks and excessive resource usage.
const maxScanDepth = 10

// loadTemplatesFromCategory loads all templates within a single category directory.
// It recursively scans subdirectories to find templates using filepath.WalkDir.
// Returns a map of templates indexed by their ID.
func loadTemplatesFromCategory(categoryPath, categoryName string) (map[string]Template, error) {
	templates := make(map[string]Template)

	err := filepath.WalkDir(categoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == categoryPath {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			log.Debug().Msgf("skipping symlink: %s", d.Name())
			return filepath.SkipDir
		}

		if !d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(categoryPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		depth := strings.Count(relPath, string(filepath.Separator)) + 1
		if depth > maxScanDepth {
			return fmt.Errorf("maximum directory depth (%d) exceeded at %s", maxScanDepth, path)
		}

		if isTemplateDirectory(path) {
			tmpl, err := LoadTemplate(path)
			if err != nil {
				return fmt.Errorf("error loading template %s: %w", d.Name(), err)
			}
			if tmpl.ID != d.Name() {
				return fmt.Errorf("template id '%s' and directory name '%s' should match", tmpl.ID, d.Name())
			}
			if existing, exists := templates[tmpl.ID]; exists {
				return fmt.Errorf("duplicate template id '%s' found (already loaded: %s)", tmpl.ID, existing.Info.Name)
			}
			templates[tmpl.ID] = tmpl
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error scanning category %s: %w", categoryName, err)
	}

	return templates, nil
}

// isTemplateDirectory checks if a directory contains an index.yaml file,
// indicating it's a template directory.
func isTemplateDirectory(dirPath string) bool {
	indexPath := filepath.Join(dirPath, "index.yaml")
	_, err := os.Stat(indexPath)
	return err == nil
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

// GetDockerComposePath finds and returns the docker-compose file path for a given template ID.
// It searches through all category directories in the templates repository to locate the template.
// Returns the absolute path to the compose file and the working directory.
func GetDockerComposePath(templateID, repoPath string) (composePath string, workingDir string, err error) {
	// Search for template in all category directories
	dirEntries, err := os.ReadDir(repoPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range dirEntries {
		if strings.HasPrefix(entry.Name(), ".") || !entry.IsDir() {
			continue
		}

		// Check if this category contains the template
		templateDir := filepath.Join(repoPath, entry.Name(), templateID)
		if !isTemplateDirectory(templateDir) {
			// Search recursively in subdirectories
			categoryPath := filepath.Join(repoPath, entry.Name())
			found, err := findTemplateInCategory(categoryPath, templateID)
			if err != nil {
				continue
			}
			if found != "" {
				templateDir = found
			} else {
				continue
			}
		}

		// Load the template to get provider config
		tmpl, err := LoadTemplate(templateDir)
		if err != nil {
			continue
		}

		if tmpl.ID != templateID {
			continue
		}

		// Get docker-compose provider config
		providerConfig, exists := tmpl.Providers["docker-compose"]
		if !exists {
			return "", "", fmt.Errorf("template %q missing docker-compose provider configuration", templateID)
		}
		if providerConfig.Path == "" {
			return "", "", fmt.Errorf("template %q docker-compose.path is empty", templateID)
		}

		// Construct the compose file path
		if filepath.IsAbs(providerConfig.Path) {
			return providerConfig.Path, filepath.Dir(providerConfig.Path), nil
		}

		composePath = filepath.Join(templateDir, providerConfig.Path)
		if _, statErr := os.Stat(composePath); statErr != nil {
			continue
		}

		return composePath, filepath.Dir(composePath), nil
	}

	return "", "", fmt.Errorf("docker-compose file for template %q not found", templateID)
}

// findTemplateInCategory recursively searches for a template directory within a category.
func findTemplateInCategory(categoryPath, templateID string) (string, error) {
	var foundPath string
	err := filepath.WalkDir(categoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			return nil
		}

		if d.Name() == templateID && isTemplateDirectory(path) {
			foundPath = path
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return foundPath, nil
}
