package zos

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// OpenFileReadOnly returns the file with read-only mode.
// Unlike [os.Open], it use [os.ModePerm] for permission.
// See also [io.OpenFile], [os.Open].
func OpenFileReadOnly(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, os.ModePerm)
}

// OpenFileWriteOnly returns the file with write-only mode.
// Unlike [os.Create], it use [os.ModePerm] for permission.
// See also [os.MkdirAll], [io.OpenFile] and [os.Create].
// It
//   - creates directory if not exists.
//   - creates file if not exists.
//   - opens file with append mode.
func OpenFileWriteOnly(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
}

// OpenFileReadWrite returns the file with read write mode.
// See also [os.MkdirAll] and [io.OpenFile].
// It
//   - creates directory if not exists.
//   - creates file if not exists.
//   - opens file with append mode.
func OpenFileReadWrite(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
}

// IsFile returns if the given path is file or not.
// It returns false for non regular files such as
//   - Directory
//   - Symbolic link
//   - Device file
//   - Unix domain socket file
//
// See [fs.ModeType].
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// IsDir returns if the given path is directory or not.
// It returns true even if the directory is symbolic link.
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// ReadFiles reads files of the given paths.
// Paths can be absolute or relative, and file or directory.
// It check all sub directories and read files if the first argument recursive is true.
// It returns an empty map and nil error if no files found.
// note that paths in the map key are cleaned by [filepath.Clean].
func ReadFiles(recursive bool, paths ...string) (map[string][]byte, error) {
	files, err := ListFiles(recursive, paths...)
	if err != nil {
		return nil, err // Return err as-is.
	}
	contents := make(map[string][]byte, len(files))
	errs := make([]error, 0, len(files))
	for _, file := range files {
		b, err := os.ReadFile(file)
		contents[filepath.Clean(file)] = b
		errs = append(errs, err)
	}
	return contents, errors.Join(errs...)
}

// ListFiles returns file paths of the given paths.
// Paths can be directories or file and can be relative path or absolute path.
// If recursive is true, it looks for all sub directories recursively.
// Note that returned paths are cleaned by [filepath.Clean].
func ListFiles(recursive bool, paths ...string) ([]string, error) {
	var files []string
	for _, path := range paths {
		if path == "" {
			continue
		}
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if !recursive {
					return fs.SkipDir
				}
				return nil
			}
			files = append(files, filepath.Clean(path))
			return nil
		})
		if err != nil && err != fs.SkipDir {
			return nil, err
		}
	}
	return files, nil
}
