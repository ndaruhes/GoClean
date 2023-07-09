package utils

import (
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadSingleFile(file []byte, fileName string, targetDir string) error {
	if err := MakeDirectory("public/" + targetDir); err != nil {
		return err
	}

	destPath := filepath.Join("public/"+targetDir, fileName)
	if err := ioutil.WriteFile(destPath, file, 0644); err != nil {
		return err
	}

	return nil
}

func DeleteFile(fileName string, targetDir string) error {
	destPath := filepath.Join("public/"+targetDir, fileName)
	fmt.Println("anjay", destPath)
	if err := os.Remove(destPath); err != nil {
		return err
	}

	return nil
}

func MakeDirectory(source string) error {
	if err := os.MkdirAll(source, 0755); err != nil {
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
