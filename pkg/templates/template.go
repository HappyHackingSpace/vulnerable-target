package templates

import (
	"fmt"
	"os"
	"path"
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
	home := path.Join(wd, "templates")
	dirEntry, err := os.ReadDir(home)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	for _, entry := range dirEntry {
		template, err := LoadTemplate(path.Join(home, entry.Name()))
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
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Author", "Technologies", "Tags"})
	for id, template := range Templates {
		if template.ID == "example-template" {
			continue
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
	}
	t.SetCaption("there are %d templates", len(Templates))
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
