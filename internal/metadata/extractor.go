package metadata

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/wtolson/go-taglib"
)

type MetaData struct {
	Title    string
	Artist   string
	Album    string
	Year     int
	Track    int
	Genre    string
	Duration time.Duration
	FilePath string
}

type Extractor struct {
	SupportedExtensions []string
}

// NewExtractor creates a new extractor instance
func NewExtractor() *Extractor {
	return &Extractor{
		SupportedExtensions: []string{".mp3", ".flac", ".mp4a", ".wav", ".ogg"},
	}
}

func (e *Extractor) isSupportedFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	for _, supportedExt := range e.SupportedExtensions {
		if supportedExt == ext {
			return true
		}

	}
	return false
}
func (e *Extractor) Extract(filePath string) (*MetaData, error) {
	if !e.isSupportedFile(filePath) {
		return nil, fmt.Errorf("unsupported file type: %s", filePath)
	}
	file, err := taglib.Read(filePath)
	if err != nil {
		return nil, fmt.Errorf("unsupported file type: %s", filePath)
	}
	defer file.Close()
	return &MetaData{
		Title:    file.Title(),
		Artist:   file.Artist(),
		Album:    file.Album(),
		Year:     file.Year(),
		Track:    file.Track(),
		Genre:    file.Genre(),
		Duration: file.Length(),
		FilePath: filePath,
	}, nil

}

func (e *Extractor) ExtractAll(dirPath string) ([]*MetaData, error) {
	var metadatalist []*MetaData

	for _, ext := range e.SupportedExtensions {
		files, err := filepath.Glob(filepath.Join(dirPath + "*" + ext))
		if err != nil {
			return nil, fmt.Errorf("Error reading directory : %s:%v", dirPath, err)
		}
		for _, file := range files {
			metadata, err := e.Extract(file)
			if err != nil {
				fmt.Printf("Warning: Could not extract metadata from %s: %v\n", file, err)
				continue
			}
			metadatalist = append(metadatalist, metadata)
		}
	}
	return metadatalist, nil
}
