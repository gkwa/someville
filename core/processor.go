package core

import (
	"github.com/go-logr/logr"
)

type FileProcessor interface {
	Process(filename, fileType string) error
}

type DirectoryWalker interface {
	Walk(basedir string, exts, ignorePaths []string, fileType string) error
}

type Processor struct {
	logger          logr.Logger
	fileProcessor   FileProcessor
	directoryWalker DirectoryWalker
}

func NewProcessor(logger logr.Logger, fileProcessor FileProcessor, directoryWalker DirectoryWalker) *Processor {
	return &Processor{
		logger:          logger,
		fileProcessor:   fileProcessor,
		directoryWalker: directoryWalker,
	}
}

func (p *Processor) ProcessFiles(basedir string, exts, ignorePaths []string, fileType string) {
	err := p.directoryWalker.Walk(basedir, exts, ignorePaths, fileType)
	if err != nil {
		p.logger.Error(err, "Error processing files")
	}
}
