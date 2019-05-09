package data

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errgo"
)

func maybeCreate(path string) {
	os.MkdirAll(path, 0700)
}

func newFolderStorage(path string) (Storage, error) {
	// create the folder if needed, make sure we have permissions
	// to write to it
	maybeCreate(path)
	fi, err := os.Stat(path)
	if err != nil {
		return nil, errgo.NoteMask(err, "no such folder: "+path)
	}
	if !fi.IsDir() {
		return nil, errgo.New("not a folder: " + path)
	}
	f, err := ioutil.TempFile(path, "")
	fileName := f.Name()
	defer f.Close()
	if err != nil {
		return nil, errgo.NoteMask(err, "could not create file in datadir")
	}
	n, err := f.Write([]byte{1, 2, 3, 4})
	if err != nil {
		return nil, errgo.NoteMask(err, "could not write to file in datadir")
	}
	if n != 4 {
		return nil, errgo.New("not all content written to file")
	}
	err = f.Close()
	if err != nil {
		return nil, errgo.NoteMask(err, "error flushing/closing file in datadir")
	}
	if err := os.Remove(fileName); err != nil {
		return nil, errgo.NoteMask(err, "could not delete file in datadir")
	}

	return folderStorage(path), nil
}

type folderStorage string

func (f folderStorage) Store(name string, data io.Reader) error {
	fi, err := os.Create(filepath.Join(string(f), name))
	defer fi.Close()
	if err != nil {
		return errgo.Mask(err)
	}
	_, err = io.Copy(fi, data)
	if err != nil {
		return errgo.Mask(err)
	}
	return errgo.Mask(fi.Close())
}
