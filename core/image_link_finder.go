package core

import (
	"regexp"
)

type ImageLinkFinder interface {
	Find(content []byte) string
}

type RegexImageLinkFinder struct {
	localImageRegex  *regexp.Regexp
	remoteImageRegex *regexp.Regexp
	anyLinkRegex     *regexp.Regexp
}

func NewRegexImageLinkFinder() *RegexImageLinkFinder {
	return &RegexImageLinkFinder{
		localImageRegex:  regexp.MustCompile(`(?i)!\[\[(.*?\.(?:png|jpg|jpeg|gif))\]\]`),
		remoteImageRegex: regexp.MustCompile(`(?i)!\[.*?\]\((https?://.*?\.(?:png|jpg|jpeg|gif))\)`),
		anyLinkRegex:     regexp.MustCompile(`(?i)\[.*?\]\((https?://[^\s)]+)\)`),
	}
}

func (rilf *RegexImageLinkFinder) Find(content []byte) string {
	if match := rilf.localImageRegex.Find(content); match != nil {
		return string(match[3 : len(match)-2])
	}

	if match := rilf.remoteImageRegex.FindSubmatch(content); len(match) > 1 {
		return string(match[1])
	}

	// If no image-specific link is found, return the first link
	if match := rilf.anyLinkRegex.FindSubmatch(content); len(match) > 1 {
		return string(match[1])
	}

	return ""
}
