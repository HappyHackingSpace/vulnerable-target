package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
)

type Template struct {
	ID        string                    `yaml:"id"`
	Info      Info                      `yaml:"info"`
	Providers map[string]ProviderConfig `yaml:"providers"`
}

type Info struct {
	Name         string                 `yaml:"name"`
	Author       string                 `yaml:"author"`
	Description  string                 `yaml:"description"`
	References   []string               `yaml:"references"`
	Technologies []string               `yaml:"technologies"`
	Tags         []string               `yaml:"tags"`
	Metadata     map[string]interface{} `yaml:"metadata"`
}

type ProviderConfig struct {
	Path string `yaml:"path"`
}

var Templates = make(map[string]Template)

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
			log.Fatal().Msgf("%v", err)
		}
		if template.ID != entry.Name() {
			log.Fatal().Msgf("id and directory name should match")
		}
		Templates[template.ID] = template
	}
}

func List() {
	ListWithFilter("")
}

// ListWithFilter lists templates filtered by tag
func ListWithFilter(filterTag string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Author", "Technologies", "Tags"})

	count := 0
	for id, template := range Templates {
		if template.ID == "example-template" {
			continue
		}

		// Apply tag filter if specified
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

		technologies := strings.Join(template.Info.Technologies, ", ")
		tags := strings.Join(template.Info.Tags, ", ")
		t.AppendRow(table.Row{
			id,
			template.Info.Name,
			template.Info.Author,
			technologies,
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

func GetByID(templateID string) (*Template, error) {
	template := Templates[templateID]
	if template.ID == "" {
		return nil, fmt.Errorf("template %s not found", templateID)
	}
	return &template, nil
}
