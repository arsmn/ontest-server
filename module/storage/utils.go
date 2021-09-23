package storage

import (
	"image"
	"image/png"
	"io"
	"os"
	"path"

	"github.com/spf13/afero"
)

func WriteFile(fs afero.Fs, dst string, in io.Reader) (os.FileInfo, error) {
	dir, _ := path.Split(dst)
	err := fs.MkdirAll(dir, 0775) //nolint:gomnd
	if err != nil {
		return nil, err
	}

	file, err := fs.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775) //nolint:gomnd
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, in)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func WriteImage(fs afero.Fs, dst string, img image.Image) (os.FileInfo, error) {
	dir, _ := path.Split(dst)
	err := fs.MkdirAll(dir, 0775) //nolint:gomnd
	if err != nil {
		return nil, err
	}

	file, err := fs.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775) //nolint:gomnd
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return info, nil
}
