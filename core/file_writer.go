package core

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

type FileWriter interface {
	UpdateContent(frontmatter map[string]interface{}, body []byte) []byte
	WriteFile(filename string, content []byte) error
}

type MarkdownFileWriter struct{}

func NewMarkdownFileWriter() *MarkdownFileWriter {
	return &MarkdownFileWriter{}
}

func (mfw *MarkdownFileWriter) UpdateContent(frontmatter map[string]interface{}, body []byte) []byte {
	frontmatterBytes, _ := yaml.Marshal(frontmatter)
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(frontmatterBytes)
	buf.WriteString("---\n")
	buf.Write(body)
	return buf.Bytes()
}

func (mfw *MarkdownFileWriter) WriteFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0o644)
}
