package web_fs_repository

import (
	"errors"
	"fmt"
	"os"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file: %s: %w", filePath, err)
		}

		return nil, fmt.Errorf("get file: %s: %w", filePath, err)
	}

	return file, nil
}
