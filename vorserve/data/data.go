package data

import (
	"io"
	"strings"

	"github.com/juju/errgo"
)

type Storage interface {
	Store(path string, data io.Reader) error
}

func NewStorage(path string) (Storage, error) {
	if strings.Contains(path, "https://") {
		return nil, errgo.New("s3 storage not yet implemented")
	}
	return newFolderStorage(path)
}
