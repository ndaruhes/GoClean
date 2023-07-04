package utils

import (
	"errors"
	"github.com/h2non/bimg"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

//func UploadSingleFile(ctx *gin.Context, key string, targetDir string) (string, error) {
//	var file, header, err = ctx.Request.FormFile(key)
//
//	fileExt := filepath.Ext(header.Filename)
//	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
//	now := time.Now().UTC()
//	fileName := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
//
//	if err := MakeDirectory("public/" + targetDir); err != nil {
//		return "", err
//	}
//
//	out, err := os.Create("public/" + targetDir + fileName)
//	if err != nil {
//		return "", err
//	}
//	defer func(out *os.File) {
//		err := out.Close()
//		if err != nil {
//			return
//		}
//	}(out)
//
//	_, err = io.Copy(out, file)
//	if err != nil {
//		return "", err
//	}
//
//	return fileName, nil
//}

func UploadSingleFile(file []byte, fileName string, targetDir string) error {
	if err := MakeDirectory("public/" + targetDir); err != nil {
		return err
	}

	destPath := filepath.Join("public/"+targetDir, fileName)
	err := ioutil.WriteFile(destPath, file, 0644)
	if err != nil {
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
