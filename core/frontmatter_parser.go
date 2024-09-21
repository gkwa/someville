package core

import (
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type FrontmatterParser interface {
	Extract(content []byte) (map[string]interface{}, []byte, error)
}

type GoldmarkFrontmatterParser struct{}

func NewGoldmarkFrontmatterParser() *GoldmarkFrontmatterParser {
	return &GoldmarkFrontmatterParser{}
}

func (gfp *GoldmarkFrontmatterParser) Extract(content []byte) (map[string]interface{}, []byte, error) {
	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	context := parser.NewContext()
	markdown.Parser().Parse(text.NewReader(content), parser.WithContext(context))

	frontmatter := meta.Get(context)
	if frontmatter == nil {
		frontmatter = make(map[string]interface{})
	}

	// Check if the content starts with "---" to determine if there's frontmatter
	lines := strings.Split(string(content), "\n")
	var body []byte
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		endIndex := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				endIndex = i
				break
			}
		}
		if endIndex != -1 {
			body = []byte(strings.Join(lines[endIndex+1:], "\n"))
		} else {
			body = content
		}
	} else {
		body = content
	}

	return frontmatter, body, nil
}
