package core

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
)

type MarkdownFileProcessor struct {
	logger            logr.Logger
	frontmatterParser FrontmatterParser
	imageLinkFinder   ImageLinkFinder
	fileWriter        FileWriter
}

func NewMarkdownFileProcessor(logger logr.Logger, frontmatterParser FrontmatterParser, imageLinkFinder ImageLinkFinder, fileWriter FileWriter) *MarkdownFileProcessor {
	return &MarkdownFileProcessor{
		logger:            logger,
		frontmatterParser: frontmatterParser,
		imageLinkFinder:   imageLinkFinder,
		fileWriter:        fileWriter,
	}
}

func (mfp *MarkdownFileProcessor) Process(filename, fileType string) error {
	mfp.logger.V(1).Info("Processing file", "file", filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	frontmatter, body, err := mfp.frontmatterParser.Extract(content)
	if err != nil {
		return fmt.Errorf("error extracting frontmatter: %w", err)
	}

	if frontmatter == nil {
		frontmatter = make(map[string]interface{})
	}

	// Check if the file has the correct filetype
	if ft, ok := frontmatter["filetype"]; !ok || ft != fileType {
		mfp.logger.V(1).Info("Skipping file with incorrect filetype", "file", filename, "filetype", ft)
		return nil
	}

	currentPic, hasPic := frontmatter["pic"].(string)
	if !hasPic || currentPic == defaultImageURL {
		imageLink := mfp.imageLinkFinder.Find(body)
		if imageLink != "" {
			frontmatter["pic"] = imageLink
			mfp.logger.V(1).Info("Updated pic in frontmatter", "file", filename, "pic", imageLink)
		} else if !hasPic {
			frontmatter["pic"] = defaultImageURL
			mfp.logger.V(1).Info("Set default pic in frontmatter", "file", filename, "pic", defaultImageURL)
		}
	}

	updatedContent := mfp.fileWriter.UpdateContent(frontmatter, body)

	if err := mfp.fileWriter.WriteFile(filename, updatedContent); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	mfp.logger.V(1).Info("File processed successfully", "file", filename)
	return nil
}
