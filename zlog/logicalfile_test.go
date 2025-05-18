package zlog

import (
	"errors"
	"os"
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestNewLogicalFile(t *testing.T) {
	t.Parallel()
	t.Run("manger fails", func(t *testing.T) {
		_, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				Pattern: "invalid.%x.%y.%z.log",
			},
		})
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("success initialize", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
		})
		defer f.Close()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		info, err := os.Stat(dir + "/application.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "files is not regular file", true, info.Mode().IsRegular())
	})
}

func TestLogicalFile_Write(t *testing.T) {
	t.Parallel()
	t.Run("no rotate", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			RotateBytes: 0,
			FileName:    "test.log",
		})
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		defer f.Close()
		f.Write([]byte("line1\n"))
		f.Write([]byte("line2\n"))
		b, err := os.ReadFile(dir + "/test.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "content not match", "line1\nline2\n", string(b))
	})
	t.Run("rotate", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			RotateBytes: 5,
			FileName:    "test.log",
		})
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		defer f.Close()
		f.Write([]byte("12345")) // No rotate.
		f.Write([]byte("67890")) // Rotate before write.
		f.Write([]byte("abc"))   // No rotate.
		f.Write([]byte("def"))   // Rotate before write.
		f.Close()
		b, err := os.ReadFile(dir + "/test.1.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "content not match", "12345", string(b))
		b, err = os.ReadFile(dir + "/test.2.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "content not match", "67890", string(b))
		b, err = os.ReadFile(dir + "/test.3.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "content not match", "abc", string(b))
	})
}

func TestLogicalFile_Close(t *testing.T) {
	t.Parallel()
	t.Run("close error", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			FileName: "test.log",
		})
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		f.curFile.Close() // Force close to make later call error.
		err = f.Close()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("rename error", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			FileName: "test.log",
		})
		f.manager.srcDir += "/not-exist/" // Make srcDir invalid.
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		err = f.Close()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
}

func TestLogicalFile_Swap(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			FileName: "test.log",
		})
		defer f.Close()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		f.Swap() // Swap 1st.
		info, err := os.Stat(dir + "/test.1.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "files is not regular file", true, info.Mode().IsRegular())
		f.Swap() // Swap 2nd.
		info, err = os.Stat(dir + "/test.2.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "files is not regular file", true, info.Mode().IsRegular())
		f.Close() // Swap 3rd.
		info, err = os.Stat(dir + "/test.3.log")
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "files is not regular file", true, info.Mode().IsRegular())
	})
	t.Run("close error", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			FileName: "test.log",
		})
		f.curFile.Close() // Force close to make Close() error.
		defer f.Close()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		err = f.Swap()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("rename error", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			FileName: "test.log",
		})
		f.manager.srcDir += "/not-exist/" // Make srcDir invalid.
		defer f.Close()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		err = f.Swap()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("open error", func(t *testing.T) {
		dir := t.TempDir()
		f, err := NewLogicalFile(&LogicalFileConfig{
			Manager: &FileManagerConfig{
				SrcDir:  dir,
				DstDir:  dir,
				Pattern: "test.%i.log",
			},
			FileName: "test.log",
		})
		f.filePath += "/not-exist/invalid.log" // Make filePath invalid.
		defer f.Close()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		err = f.Swap()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
}

func TestLogicalFile_Fallback(t *testing.T) {
	t.Parallel()
	t.Run("nil error", func(t *testing.T) {
		var onFallbackCalled bool
		f := &LogicalFile{
			onFallback: func(err error) {
				onFallbackCalled = true
			},
			isStderr: false,
			curSize:  100, // Set to non zero.
		}
		f.fallbackStderr(nil, "fallback test")
		ztesting.AssertEqual(t, "curFile was unexpectedly replaced", nil, f.curFile)
		ztesting.AssertEqual(t, "isStderr should be false", false, f.isStderr)
		ztesting.AssertEqual(t, "curSize should not modified", 100, f.curSize)
		ztesting.AssertEqual(t, "onFallback should be false", false, onFallbackCalled)
	})
	t.Run("non-nil error", func(t *testing.T) {
		var onFallbackCalled bool
		f := &LogicalFile{
			onFallback: func(err error) {
				onFallbackCalled = true
			},
			isStderr: false,
			curSize:  100, // Set to non zero.
		}
		f.fallbackStderr(errors.New("non-nil"), "fallback test")
		ztesting.AssertEqual(t, "curFile is not stderr", os.Stderr, f.curFile)
		ztesting.AssertEqual(t, "isStderr should be true", true, f.isStderr)
		ztesting.AssertEqual(t, "curSize should be zero", 0, f.curSize)
		ztesting.AssertEqual(t, "onFallback is not called", true, onFallbackCalled)
	})
}
