package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
	}
}

func (ls *LocalStorage) Save(src, dest string) (string, error) {
	destPath := filepath.Join(ls.BasePath, dest)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", err
	}

	// Copy file
	srcFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func (ls *LocalStorage) Get(path string) (string, error) {
	fullPath := filepath.Join(ls.BasePath, path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", path)
	}
	return fullPath, nil
}

func (ls *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(ls.BasePath, path)
	return os.Remove(fullPath)
}

func (ls *LocalStorage) List(prefix string) ([]string, error) {
	searchPath := filepath.Join(ls.BasePath, prefix)
	var files []string

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(ls.BasePath, path)
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}
