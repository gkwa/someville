package core

import (
	"github.com/go-logr/logr"
)

func NewDefaultProcessor(logger logr.Logger) *Processor {
	frontmatterParser := NewGoldmarkFrontmatterParser()
	imageLinkFinder := NewRegexImageLinkFinder()
	fileWriter := NewMarkdownFileWriter()
	fileProcessor := NewMarkdownFileProcessor(logger, frontmatterParser, imageLinkFinder, fileWriter)
	directoryWalker := NewFileSystemWalker(logger, fileProcessor)
	return NewProcessor(logger, fileProcessor, directoryWalker)
}
