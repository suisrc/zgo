package files

import (
	"errors"
	"os"
	"path/filepath"
)

// GetFile 获取文件
func GetFile(path string) (*os.File, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pathX := path
	if !filepath.IsAbs(path) {
		pathX = filepath.Join(pwd, path)
	}

	// 检查组合路径文件是否存在(relative path)
	if _, err := os.Stat(pathX); !os.IsNotExist(err) {
		return os.Open(pathX)
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, errors.New("file is not exist")
	}
	return f, err
}
