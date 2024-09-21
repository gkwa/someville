package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
)

type FileSystemWalker struct {
	logger        logr.Logger
	fileProcessor FileProcessor
}

func NewFileSystemWalker(logger logr.Logger, fileProcessor FileProcessor) *FileSystemWalker {
	return &FileSystemWalker{
		logger:        logger,
		fileProcessor: fileProcessor,
	}
}

func (fsw *FileSystemWalker) Walk(basedir string, exts, ignorePaths []string, fileType string) error {
	return filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fsw.logger.Error(err, "Error accessing path", "path", path)
			return nil
		}

		if info.IsDir() {
			if fsw.shouldIgnoreDir(path, ignorePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		for _, ext := range exts {
			if strings.EqualFold(filepath.Ext(path), "."+ext) {
				if err := fsw.fileProcessor.Process(path, fileType); err != nil {
					fsw.logger.Error(err, "Error processing file", "file", path)
				}
				break
			}
		}

		return nil
	})
}

func (fsw *FileSystemWalker) shouldIgnoreDir(path string, ignorePaths []string) bool {
	lowercasePath := strings.ToLower(path)
	for _, ignorePath := range ignorePaths {
		if strings.Contains(lowercasePath, strings.ToLower(ignorePath)) {
			return true
		}
	}
	return false
}
