package validators

import (
	"bytes"
	"errors"
	"fmt"
	"go-clean/src/models/messages"
	"go-clean/src/models/validations"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"golang.org/x/exp/slices"
)

var imageValidation = &validations.ImageValidation{
	MaxSize:   2,
	MinWidth:  300,
	MaxWidth:  640,
	MinHeight: 300,
	MaxHeight: 640,
	Format:    []string{"jpeg", "png"},
}

func ValidateImage(header *multipart.FileHeader) error {
	file, err := header.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if imageValidation.MaxSize > 0 && float32(header.Size) > imageValidation.MaxSize*1024*1024 {
		return errors.New(fmt.Sprintf("file size exceeds %f megabytes", imageValidation.MaxSize))
	}
	img, format, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	if !slices.Contains(imageValidation.Format, format) {
		return errors.New("invalid image format")
	}

	if imageValidation.MinWidth > 0 && img.Bounds().Dx() < imageValidation.MinWidth || imageValidation.MinHeight > 0 && img.Bounds().Dy() < imageValidation.MinHeight || imageValidation.MaxWidth > 0 && img.Bounds().Dx() > imageValidation.MaxWidth || imageValidation.MaxHeight > 0 && img.Bounds().Dy() > imageValidation.MaxHeight {
		//return errors.New(fmt.Sprintf("image size must be between %dx%d and %dx%d", imageValidation.MinWidth, imageValidation.MinHeight, imageValidation.MaxWidth, imageValidation.MaxHeight))
		//ErrorCode:  fmt.Sprintf("image size must be between %dx%d and %dx%d", imageValidation.MinWidth, imageValidation.MinHeight, imageValidation.MaxWidth, imageValidation.MaxHeight),

		return &messages.ErrorWrapper{
			ErrorCode:  "ERROR-FILE-50001",
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}
