package articles

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed data/articles.yaml
var Fixture embed.FS

type Article struct {
	ID              string         `yaml:"id" json:"id"`
	Category        string         `yaml:"category" json:"category"`
	Title           string         `yaml:"title" json:"title"`
	Summary         string         `yaml:"summary" json:"summary"`
	Audience        string         `yaml:"audience" json:"audience"`
	Materials       []string       `yaml:"materials" json:"materials"`
	Steps           []string       `yaml:"steps" json:"steps"`
	Troubleshooting []Troubleshoot `yaml:"troubleshooting" json:"troubleshooting"`
}

type Troubleshoot struct {
	Symptom  string `yaml:"symptom" json:"symptom"`
	Solution string `yaml:"solution" json:"solution"`
}

type fixture struct {
	Articles []Article `yaml:"articles"`
}

func LoadYAML(fsys embed.FS) ([]Article, error) {
	b, err := fsys.ReadFile("data/articles.yaml")
	if err != nil {
		return nil, fmt.Errorf("read article fixture: %w", err)
	}
	var value fixture
	if err := yaml.Unmarshal(b, &value); err != nil {
		return nil, fmt.Errorf("decode article fixture: %w", err)
	}
	if len(value.Articles) == 0 {
		return nil, fmt.Errorf("article fixture is empty")
	}
	return value.Articles, nil
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
