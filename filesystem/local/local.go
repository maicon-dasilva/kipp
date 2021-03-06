package local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/uhthomas/kipp/filesystem"
)

// A FileSystem contains information about the local filesystem.
type FileSystem struct{ dir, tmp string }

// New creates a new FileSystem, and makes the relevant directories for
// dir and tmp.
func New(dir string) (*FileSystem, error) {
	tmp := filepath.Join(dir, "tmp")
	if err := os.MkdirAll(tmp, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &FileSystem{dir: dir, tmp: tmp}, nil
}

// Create creates a temporary file, and wraps it so when the file is closed,
// it is moved to name.
func (fs FileSystem) Create(_ context.Context, name string) (filesystem.Writer, error) {
	f, err := ioutil.TempFile(fs.tmp, "kipp")
	if err != nil {
		return nil, fmt.Errorf("temp file: %w", err)
	}
	return &writer{f: f, name: filepath.Join(fs.dir, name)}, nil
}

// Open opens the named file.
func (fs FileSystem) Open(_ context.Context, name string) (filesystem.Reader, error) {
	return os.Open(filepath.Join(fs.dir, name))
}

// Remove removes the named file.
func (fs FileSystem) Remove(_ context.Context, name string) error {
	return os.Remove(name)
}
