package storage

import (
	"os"
)

type FileStorage struct {
	file *os.File
}

func (fs *FileStorage) Set(key, value string) error {

}

func (fs *FileStorage) Get(key string) (string, error) {

}
