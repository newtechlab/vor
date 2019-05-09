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
	parts := strings.SplitN(path, ":", 2)
	if len(parts) != 2 {
		return nil, errgo.New("unrecognized data specifier, should start with s3: or file:")
	}
	typ, arg := parts[0], parts[1]
	if arg == "" {
		return nil, errgo.New("bad data specifier")
	}

	switch typ {
	case "s3":
		return newS3Storage(arg)
	case "file":
		return newFolderStorage(arg)
	}
	return nil, errgo.New("unknown data type specifier")
}
