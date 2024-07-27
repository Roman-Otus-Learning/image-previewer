package filesystem

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

var _ Filesystem = (*DiscFilesystem)(nil)

var ErrFileNotExists = errors.New("file is not exists")

const filePermission = 0o700

type Filesystem interface {
	WriteFile(name string, content []byte) error
	ReadFile(name string) ([]byte, error)
	RemoveFile(name string) error
}

type DiscFilesystem struct {
	basePath string
}

func CreateDiskFilesystem(basePath string) (*DiscFilesystem, error) {
	if err := ensureDir(basePath); err != nil {
		return nil, errors.Wrap(err, "ensure file dir exists")
	}

	return &DiscFilesystem{
		basePath: basePath,
	}, nil
}

func (f *DiscFilesystem) WriteFile(name string, content []byte) error {
	path := filepath.Join(f.basePath, name)

	if err := os.WriteFile(path, content, filePermission); err != nil {
		return errors.Wrap(err, "filesystem write file")
	}

	return nil
}

func (f *DiscFilesystem) ReadFile(name string) ([]byte, error) {
	path := filepath.Join(f.basePath, name)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExists
		}

		return nil, errors.Wrap(err, "filesystem read file")
	}

	return content, nil
}

func (f *DiscFilesystem) RemoveFile(name string) error {
	path := filepath.Join(f.basePath, name)

	if err := os.Remove(path); err != nil {
		return errors.Wrap(err, "remove file")
	}

	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(dir, filePermission); err != nil {
		return errors.Wrap(err, "ensure dir")
	}

	return nil
}
