package utils

import (
	"mime"
	"mime/multipart"
	"path/filepath"
)

func FileHeadertContentType(fileHeader *multipart.FileHeader) string {
	// Try to get from header first
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "" {
		return contentType
	}

	// Fallback to detection by extension
	ext := filepath.Ext(fileHeader.Filename)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	// Default fallback
	return "application/octet-stream"
}
