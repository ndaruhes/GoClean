package services

import (
	"bytes"
	"context"
	"fmt"
	"go-clean/src/models/messages"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

type GoogleStorageConfig struct {
	Path     string
	Filename string
	Bucket   string
}

type StorageService struct{}

func NewStorageService() *StorageService {
	return &StorageService{}
}

func (service *StorageService) UploadFromBytes(ctx *gin.Context, b []byte, config GoogleStorageConfig) error {
	parentCtx := ctx.Request.Context() // Extract the underlying context

	timeout := time.Now().Add(time.Second * 50)
	ctxWithTimeout, cancel := context.WithDeadline(parentCtx, timeout)
	defer cancel()

	reader := bytes.NewReader(b)

	client, err := storage.NewClient(ctxWithTimeout)
	if messages.HasError(err) {
		return err
	}
	defer client.Close()

	wc := client.Bucket(config.Bucket).Object(fmt.Sprintf("%s/%s", config.Path, config.Filename)).NewWriter(ctxWithTimeout)

	if _, err := io.Copy(wc, reader); messages.HasError(err) {
		return err
	}
	if err := wc.Close(); messages.HasError(err) {
		return err
	}

	return nil
}
