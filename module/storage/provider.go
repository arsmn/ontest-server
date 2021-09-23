package storage

import "github.com/spf13/afero"

type Provider interface {
	FS() afero.Fs
}
