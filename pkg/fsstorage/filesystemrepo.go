package fsstorage

import (
	"bufio"
	"os"
)

type FileSystemRepo struct{}

func NewFileSystemRepo() *FileSystemRepo {
	return &FileSystemRepo{}
}

func (repo *FileSystemRepo) Save(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() error {
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}()
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(content); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (repo *FileSystemRepo) Read(filename string) (string, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(f), nil
}
