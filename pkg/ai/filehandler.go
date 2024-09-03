package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (fh *FileHandler) ProcessPaths(paths []string) (string, error) {
	var context strings.Builder

	for _, path := range paths {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			context.WriteString(fmt.Sprintf("// %s\n%s\n\n", filePath, string(content)))
			return nil
		})

		if err != nil {
			return "", err
		}
	}

	return context.String(), nil
}

func (fh *FileHandler) LoadPrompt(promptFile string) (string, error) {
	content, err := os.ReadFile(promptFile)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
