package utils

import (
	"bytes"
	"fmt"
	"image"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/exceptions"
)

type ImageValidator struct {
	MaxFileSize   int64
	MaxDimensions int
	AllowedTypes  []string
}

// ValidateBasicProperties checks file size and basic mime type
func (v *ImageValidator) ValidateBasicProperties(file *multipart.FileHeader) (bool, exceptions.PictureException) {
	if file == nil {
		return false, exceptions.NewPictureException(exceptions.NoFileProvided)
	}

	// Check file size
	if file.Size > v.MaxFileSize {
		return false, exceptions.NewPictureException(exceptions.ExceedsMaxSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !v.isAllowedExtension(ext) {
		return false, exceptions.NewPictureException(exceptions.ExtensionNotAllowed)
	}

	return true, exceptions.NewPictureException(exceptions.NoException)
}

// ValidateFileType performs deep file type checking
func (v *ImageValidator) ValidateFileType(fileBytes []byte) (*types.Type, exceptions.PictureException) {
	// Check if file is an image using filetype library
	kind, err := filetype.Match(fileBytes)
	if err != nil {
		return nil, exceptions.NewPictureException(exceptions.ErrServerError)
	}

	if kind == filetype.Unknown {
		return nil, exceptions.NewPictureException(exceptions.UnknownExtension)
	}

	if !v.isAllowedMimeType(kind.MIME.Value) {
		fmt.Println("mime: ", kind.MIME.Value)
		return nil, exceptions.NewPictureException(exceptions.MimeTypeNotAllowed)
	}

	return &kind, exceptions.NewPictureException(exceptions.NoException)
}

// ValidateImageIntegrity checks image dimensions and ensures it's a valid image
func (v *ImageValidator) ValidateImageIntegrity(fileBytes []byte) (image.Image, exceptions.PictureException) {
	// Decode image to verify integrity
	img, format, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, exceptions.NewPictureException(exceptions.ErrServerError)

	}

	// Check image dimensions
	bounds := img.Bounds()
	if bounds.Dx() > v.MaxDimensions || bounds.Dy() > v.MaxDimensions {
		// fmt.Errorf("image dimensions exceed maximum allowed size of %dx", v.MaxDimensions)
		return nil, exceptions.NewPictureException(exceptions.ExceedsMaxSize)
	}

	// Additional format-specific validation could be added here
	switch format {
	case "jpeg", "png", "gif":
		// Format is already validated, but you could add format-specific checks here
	default:
		return nil, exceptions.NewPictureException(exceptions.UnknownExtension)
	}

	return img, exceptions.NewPictureException(exceptions.NoException)
}

func (v *ImageValidator) isAllowedExtension(ext string) bool {
	for _, allowed := range v.AllowedTypes {
		if "image/"+strings.TrimLeft(ext, ".") == allowed {
			return true
		}
		if "."+allowed == ext {
			return true
		}
	}
	return false
}

func (v *ImageValidator) isAllowedMimeType(mime string) bool {
	for _, allowed := range v.AllowedTypes {
		if strings.Contains(mime, allowed) {
			return true
		}
	}
	return false
}
