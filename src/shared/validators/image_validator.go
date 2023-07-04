package validators

import (
	"errors"
	"fmt"
	"go-clean/models/validations"
	"golang.org/x/exp/slices"
	"image"
	"mime/multipart"
)

func ValidateImage(header *multipart.FileHeader, validation validations.ImageValidation) error {
	file, err := header.Open()
	if err != nil {
		return err
	}

	if validation.MaxSize > 0 && float32(header.Size) > validation.MaxSize*1024*1024 {
		return errors.New(fmt.Sprintf("file size exceeds %f megabytes", validation.MaxSize))
	}
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	if !slices.Contains(validation.Format, format) {
		return errors.New("invalid image format")
	}

	if validation.MinWidth > 0 && img.Bounds().Dx() < validation.MinWidth || validation.MinHeight > 0 && img.Bounds().Dy() < validation.MinHeight || validation.MaxWidth > 0 && img.Bounds().Dx() > validation.MaxWidth || validation.MaxHeight > 0 && img.Bounds().Dy() > validation.MaxHeight {
		return errors.New(fmt.Sprintf("image size must be between %dx%d and %dx%d", validation.MinWidth, validation.MinHeight, validation.MaxWidth, validation.MaxHeight))
	}
	return nil
}
