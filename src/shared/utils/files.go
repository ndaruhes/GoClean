package utils

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/h2non/bimg"
)

func MakeDirectory(source string) error {
	if err := os.MkdirAll(source, 0755); err != nil {
		return err
	}

	return nil
}

func FileExists(targetDir string, fileName string) error {
	destDir := filepath.Join(targetDir, fileName)
	_, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	return nil
}

func UploadSingleFile(file []byte, targetDir string, fileName string) error {
	if err := MakeDirectory(targetDir); err != nil {
		return err
	}

	destDir := filepath.Join(targetDir, fileName)
	if err := ioutil.WriteFile(destDir, file, 0644); err != nil {
		return err
	}

	return nil
}

func MoveSingleFile(sourceDir string, targetDir string, fileName string) error {
	sourcePath := filepath.Join(sourceDir, fileName)
	targetPath := filepath.Join(targetDir, fileName)

	if err := MakeDirectory(targetDir); err != nil {
		return err
	}

	err := os.Rename(sourcePath, targetPath)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSingleFile(targetDir string, fileName string) error {
	if err := FileExists(targetDir, fileName); err != nil {
		return err
	}

	destDir := filepath.Join(targetDir, fileName)
	if err := os.Remove(destDir); err != nil {
		return err
	}

	return nil
}

func MultipartFileHeaderToByte(header *multipart.FileHeader) ([]byte, error) {
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

func CompressFile(file []byte, quality int) ([]byte, error) {
	if quality < 0 || quality > 100 {
		return nil, errors.New("compression quality must be between 0 and 100")
	}

	converted, err := bimg.NewImage(file).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
}
