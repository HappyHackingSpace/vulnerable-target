package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
)

// Template represents a vulnerable target environment configuration.
type Template struct {
	ID             string                    `yaml:"id"`
	Info           Info                      `yaml:"info"`
	ProofOfConcept map[string][]string       `yaml:"poc"`
	Remediation    []string                  `yaml:"remediation"`
	Providers      map[string]ProviderConfig `yaml:"providers"`
	PostInstall    []string                `yaml:"post-install"`
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

// Templates contains all loaded templates indexed by their ID.
var Templates = make(map[string]Template)

// Init loads all templates from the templates directory.
func Init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	home := filepath.Join(wd, "templates")
	dirEntry, err := os.ReadDir(home)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	for _, entry := range dirEntry {
		template, err := LoadTemplate(filepath.Join(home, entry.Name()))
		if err != nil {
			log.Fatal().Msgf("Error loading template %s: %v", entry.Name(), err)
		}
		if template.ID != entry.Name() {
			log.Fatal().Msgf("id and directory name should match")
		}
		Templates[template.ID] = template
	}
}

// List displays all available templates in a table format.
func List() {
	ListWithFilter("")
}

// ListWithFilter displays all available templates in a table format, optionally filtered by tag.
func ListWithFilter(filterTag string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Author", "Targets", "Type", "Tags"})

	count := 0
	for _, template := range Templates {
		if filterTag != "" {
			hasTag := false
			for _, tag := range template.Info.Tags {
				if strings.EqualFold(tag, filterTag) || strings.Contains(strings.ToLower(tag), strings.ToLower(filterTag)) {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		tags := strings.Join(template.Info.Tags, ", ")
		targets := strings.Join(template.Info.Targets, ", ")
		t.AppendRow(table.Row{
			template.ID,
			template.Info.Name,
			template.Info.Author,
			targets,
			template.Info.Type,
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

// GetByID retrieves a template by its ID.
func GetByID(templateID string) (*Template, error) {
	template := Templates[templateID]
	if template.ID == "" {
		return nil, fmt.Errorf("template %s not found", templateID)
	}
	return &template, nil
}