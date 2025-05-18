package zlog

import (
	"cmp"
	"os"
	"path/filepath"
	"sync"

	"github.com/aileron-projects/go/zruntime"
)

// LogicalFileConfig is the configuration for [LogicalFile].
// Use [NewLogicalFile] to create a new instance of it.
type LogicalFileConfig struct {
	// Manager is the archive files manager config.
	Manager *FileManagerConfig
	// RotateBytes is the maximum physical file size in bytes.
	// Physical file will be rotated when reached to the size.
	// Zero or negative means no limit, or no file rotation.
	RotateBytes int64
	// FileName is the physical file name that is
	// actively written to.
	// If empty, "application.log" is used.
	FileName string
	// OnFallback is called when underlying physical file
	// operation failed. Before calling the function,
	// write target will be replaced to [os.Stderr].
	// If not set, [os.Stderr] will be used until
	// next file rotation occur.
	OnFallback func(error)
}

// NewLogicalFile returns a new instance of [LogicalFile].
func NewLogicalFile(c *LogicalFileConfig) (*LogicalFile, error) {
	m, err := NewFileManager(c.Manager)
	if err != nil {
		return nil, err
	}
	c.FileName = cmp.Or(c.FileName, "application.log")
	f := &LogicalFile{
		rotBytes:   c.RotateBytes,
		filePath:   filepath.Join(m.srcDir, c.FileName),
		manager:    m,
		onFallback: c.OnFallback,
	}
	return f, f.Swap()
}

// LogicalFile is a logical file type.
// It implements [io.Writer] interface.
// Users do not need to care physical file management
// such as open, close, rename or remove.
// Use [NewLogicalFile] to create a new instance of [LogicalFile].
type LogicalFile struct {
	mu sync.Mutex
	// manager manages archived files.
	manager *FileManager
	// rotBytes is the maximum file size in byte.
	// Current active file will be swapped to a new one
	// when the current file size exceeded the size.
	// Zero or negative means no limit, or no file rotation.
	rotBytes int64
	// curSize is the current file size in byte.
	curSize int64
	// filePath is the active file path.
	// This should contain directories.
	filePath string
	// curFile is the current physical file.
	// [os.Stderr] can be used as fallback when some file operation failed.
	curFile *os.File
	// isStderr is the flag if the [os.Stderr] is used as the curFile.
	isStderr bool
	// onFallback is the hook function that will be called
	// when physical file operation failed.
	// onFallback will be called after fallback, from physical file to stderr,
	// had been completed with non-nil error.
	onFallback func(error)
}

// Write writes the given data in to the file.
// It will be immediately written to the underlying physical file.
// Write implements [io.Writer.Write].
// Write is safe for concurrent call.
func (f *LogicalFile) Write(b []byte) (int, error) {
	f.mu.Lock() // Protect f.curFile. It should not be swapped now.
	defer f.mu.Unlock()
	if f.rotBytes > 0 && f.curSize+int64(len(b)) > f.rotBytes {
		_ = f.swap() // Ignore swap error. Write should not depends on the error.
	}
	n, err := f.curFile.Write(b)
	f.curSize += int64(n)
	f.fallbackStderr(err, "failed to write file.") // Fallback if non nil.
	return n, err
}

// Close closes the file.
// It closes the underlying physical file.
// Unlike [LogicalFile.Swap]. it does not create a new file.
// Close implements [io.Closer.Close].
// Close is safe for concurrent call.
func (f *LogicalFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.close()
}

func (f *LogicalFile) close() error {
	if f.curFile == nil || f.isStderr {
		f.curSize = 0      // Reset.
		f.curFile = nil    // Reset.
		f.isStderr = false // Reset.
		return nil
	}
	name := f.curFile.Name() // Keep name.
	if err := f.curFile.Close(); err != nil {
		return err
	}
	f.curSize = 0   // Reset.
	f.curFile = nil // Reset.
	if err := os.Rename(name, f.manager.NewFile()); err != nil {
		return err
	}
	return f.manager.Manage()
}

// Swap swaps the active file to new one.
// [FileManager.Manage] will be called internally in a new goroutine
// which means archived files are managed asynchronously.
// Swap is safe for concurrent call.
func (f *LogicalFile) Swap() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.swap()
}

func (f *LogicalFile) swap() error {
	if f.curFile != nil && !f.isStderr {
		if err := f.curFile.Close(); err != nil {
			f.fallbackStderr(err, "failed to close file.")
			return err
		}
		if err := os.Rename(f.curFile.Name(), f.manager.NewFile()); err != nil {
			f.fallbackStderr(err, "failed to rename file.")
			return err
		}
	}

	go func() {
		err := f.manager.Manage()                              // Manage archives in parallel.
		zruntime.ReportErr(err, "zlog: error managing files.") // Report only if err is non-nil.
	}()

	ff, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		f.fallbackStderr(err, "failed to open file.")
		return err
	}
	f.curSize = 0
	f.isStderr = false
	f.curFile = ff
	if info, err := ff.Stat(); err == nil {
		f.curSize = info.Size() // Update size when available.
	}
	return nil
}

// fallbackStderr replaces f.curFile to [os.Stderr].
// This should be called when physical file operation failed.
func (f *LogicalFile) fallbackStderr(err error, msg string) {
	if err == nil {
		return
	}
	zruntime.ReportErr(err, "zlog: "+msg+" fallback to stderr.") // Report if error.
	f.curFile = os.Stderr                                        // Fallback until next swap.
	f.isStderr = true
	f.curSize = 0
	if f.onFallback != nil {
		f.onFallback(err)
	}
}
