package service

import (
	"awesomeProject/utils"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileUploader struct {
	TargetDir    string
	AllowedExits []string
}

var targetDir = utils.GoDotEnvVariable("FILE_UPLOAD_PATH")

func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size == 0 {
		return "", fmt.Errorf("file is required")
	}

	fu := FileUploader{
		TargetDir:    targetDir,
		AllowedExits: []string{".jpg", ".jpeg", ".png", ".pdf", ".docx"},
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := false
	for _, allowedExt := range fu.AllowedExits {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("invalid file type")
	}

	fileName := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(fileHeader.Filename))
	targetFilePath := filepath.Join(fu.TargetDir, fileName)

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
