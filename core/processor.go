package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

const defaultImageURL = "https://encrypted-tbn2.gstatic.com/shopping?q=tbn:ANd9GcS71yfHYfcINhCdWC_V6hy6tSK-UqwHi2GbC1TKTXHRGsoJEuyC9rfjc11Nj6J2jIcqf07qnF6Lpp973qeWH8j5I2GCtHpd1rRBXegWkCIm4AcYDGGCAkEmfQ"

func ProcessFiles(logger logr.Logger, basedir string, exts, ignorePaths []string, fileType string) {
	for _, ext := range exts {
		err := filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Error(err, "Error accessing path", "path", path)
				return nil
			}

			if info.IsDir() {
				if shouldIgnoreDir(path, ignorePaths) {
					return filepath.SkipDir
				}
				return nil
			}

			if filepath.Ext(path) == "."+ext {
				if err := processFile(logger, path, fileType); err != nil {
					logger.Error(err, "Error processing file", "file", path)
				}
			}

			return nil
		})
		if err != nil {
			logger.Error(err, "Error walking directory", "basedir", basedir)
		}
	}
}

func shouldIgnoreDir(path string, ignorePaths []string) bool {
	lowercasePath := strings.ToLower(path)
	for _, ignorePath := range ignorePaths {
		if strings.Contains(lowercasePath, strings.ToLower(ignorePath)) {
			return true
		}
	}
	return false
}

func processFile(logger logr.Logger, filename, fileType string) error {
	logger.V(1).Info("Processing file", "file", filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	frontmatter, body, err := extractFrontmatter(content)
	if err != nil {
		return fmt.Errorf("error extracting frontmatter: %w", err)
	}

	if frontmatter == nil {
		frontmatter = make(map[string]interface{})
	}

	// Check if the file has the correct filetype
	if ft, ok := frontmatter["filetype"]; !ok || ft != fileType {
		logger.V(1).Info("Skipping file with incorrect filetype", "file", filename, "filetype", ft)
		return nil
	}

	if _, ok := frontmatter["pic"]; !ok {
		imageLink := findFirstImageLink(body)
		if imageLink == "" {
			imageLink = defaultImageURL
		}
		frontmatter["pic"] = imageLink
	}

	updatedContent := updateFileContent(frontmatter, body)

	if err := os.WriteFile(filename, updatedContent, 0o644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	logger.V(1).Info("File processed successfully", "file", filename)
	return nil
}

func extractFrontmatter(content []byte) (map[string]interface{}, []byte, error) {
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

func findFirstImageLink(content []byte) string {
	localImageRegex := regexp.MustCompile(`!\[\[(.*?\.(?:png|jpg|jpeg|gif))\]\]`)
	remoteImageRegex := regexp.MustCompile(`!\[.*?\]\((https?://.*?\.(?:png|jpg|jpeg|gif))\)`)

	if match := localImageRegex.Find(content); match != nil {
		return string(match[3 : len(match)-2])
	}

	if match := remoteImageRegex.FindSubmatch(content); len(match) > 1 {
		return string(match[1])
	}

	return ""
}

func updateFileContent(frontmatter map[string]interface{}, body []byte) []byte {
	frontmatterBytes, _ := yaml.Marshal(frontmatter)
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(frontmatterBytes)
	buf.WriteString("---\n")
	buf.Write(body)
	return buf.Bytes()
}
