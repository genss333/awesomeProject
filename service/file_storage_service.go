package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileUploader struct {
	targetDir    string
	allowedExits []string
}

func NewFileUploader(targetDir string, allowedExits []string) *FileUploader {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		err := os.MkdirAll(targetDir, 0777)
		if err != nil {
			return nil
		}
	}
	return &FileUploader{targetDir: targetDir, allowedExits: allowedExits}
}

func (fu *FileUploader) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size == 0 {
		return "", fmt.Errorf("file is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := false
	for _, allowedExt := range fu.allowedExits {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("invalid file type")
	}

	fileName := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(fileHeader.Filename))
	targetFilePath := filepath.Join(fu.targetDir, fileName)

	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func(srcFile multipart.File) {
		err := srcFile.Close()
		if err != nil {
			return
		}
	}(srcFile)

	destFile, err := os.Create(targetFilePath)
	if err != nil {
		return "", err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			return
		}
	}(destFile)

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", err
	}

	return targetFilePath, nil
}
