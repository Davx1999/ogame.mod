package v104

import (
	v10 "github.com/alaingilbert/ogame/pkg/extractor/v10"
)

// Extractor ...
type Extractor struct {
	v10.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractUpgradeToken ...
func (e *Extractor) ExtractUpgradeToken(pageHTML []byte) (string, error) {
	return extractUpgradeToken(pageHTML)
}
