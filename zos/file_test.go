package zos_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/aileron-projects/go/zos"
	"github.com/aileron-projects/go/ztesting"
)

func TestOpenFileReadOnly(t *testing.T) {
	t.Run("cannot write", func(t *testing.T) {
		f, err := zos.OpenFileReadOnly("./testdata/readonly.txt")
		defer f.Close()
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		n, err := f.Write([]byte("test"))
		ztesting.AssertEqual(t, "could write", 0, n)
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
}

func TestOpenFileWriteOnly(t *testing.T) {
	t.Parallel()
	defer func() {
		os.RemoveAll("./testdata/TestOpenFileWriteOnly/")
	}()

	t.Run("cannot read", func(t *testing.T) {
		f, err := zos.OpenFileWriteOnly("./testdata/TestOpenFileWriteOnly/cannot-read1.txt")
		defer f.Close()
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		n, err := f.Write([]byte("test"))
		ztesting.AssertEqual(t, "unexpected written bytes", 4, n)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		n, err = f.ReadAt(make([]byte, 10), 0)
		ztesting.AssertEqual(t, "could read", 0, n)
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
	t.Run("create dir failed", func(t *testing.T) {
		_, err := zos.OpenFileWriteOnly("./testdata/TestOpenFileWriteOnly\x00/cannot-read2.txt")
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
}

func TestOpenFileReadWrite(t *testing.T) {
	t.Parallel()
	defer func() {
		os.RemoveAll("./testdata/TestOpenFileReadWrite/")
	}()

	t.Run("cannot read", func(t *testing.T) {
		f, err := zos.OpenFileReadWrite("./testdata/TestOpenFileReadWrite/read-write1.txt")
		defer f.Close()
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		n, err := f.Write([]byte("test"))
		ztesting.AssertEqual(t, "unexpected written bytes", 4, n)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
		n, err = f.ReadAt(make([]byte, 3), 0)
		ztesting.AssertEqual(t, "could not read", 3, n)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("create dir failed", func(t *testing.T) {
		_, err := zos.OpenFileReadWrite("./testdata/TestOpenFileReadWrite\x00/read-write2.txt")
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
}

func TestIsFile(t *testing.T) {
	t.Parallel()
	t.Run("file", func(t *testing.T) {
		ok, err := zos.IsFile("./testdata/regular-file.txt")
		ztesting.AssertEqual(t, "incorrect judgement", true, ok)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("non file", func(t *testing.T) {
		ok, err := zos.IsFile("./testdata/")
		ztesting.AssertEqual(t, "incorrect judgement", false, ok)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("stat error", func(t *testing.T) {
		ok, err := zos.IsFile("./testdata\x00/")
		ztesting.AssertEqual(t, "incorrect judgement", false, ok)
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
}

func TestIsDir(t *testing.T) {
	t.Parallel()
	t.Run("dir", func(t *testing.T) {
		ok, err := zos.IsDir("./testdata/")
		ztesting.AssertEqual(t, "incorrect judgement", true, ok)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("non dir", func(t *testing.T) {
		ok, err := zos.IsDir("./testdata/regular-file.txt")
		ztesting.AssertEqual(t, "incorrect judgement", false, ok)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("stat error", func(t *testing.T) {
		ok, err := zos.IsDir("./testdata\x00/")
		ztesting.AssertEqual(t, "incorrect judgement", false, ok)
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
}

func TestReadFiles(t *testing.T) {
	t.Parallel()
	t.Run("no paths", func(t *testing.T) {
		contents, err := zos.ReadFiles(false)
		ztesting.AssertEqual(t, "returned data not match", 0, len(contents))
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("empty path", func(t *testing.T) {
		contents, err := zos.ReadFiles(false, "")
		ztesting.AssertEqual(t, "returned data not match", 0, len(contents))
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("invalid path", func(t *testing.T) {
		contents, err := zos.ReadFiles(false, "./not-found/")
		ztesting.AssertEqual(t, "returned data not match", 0, len(contents))
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
	t.Run("not recursive", func(t *testing.T) {
		contents, err := zos.ReadFiles(false, "./testdata/regular-file.txt", "./testdata/files/")
		want := map[string][]byte{filepath.Clean("./testdata/regular-file.txt"): []byte("regular-file")}
		for k, v := range want {
			vv := contents[k]
			if !slices.Equal(v, vv) {
				ztesting.AssertEqualSlice(t, "content not match for key="+k, v, contents[k])
			}
		}
		ztesting.AssertEqual(t, "returned data size not match", 1, len(contents))
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("recursive", func(t *testing.T) {
		contents, err := zos.ReadFiles(true, "./testdata/files/")
		want := map[string][]byte{}
		want[filepath.Clean("./testdata/files/level0.txt")] = []byte("level0")
		want[filepath.Clean("./testdata/files/level1/level1.txt")] = []byte("level1")
		for k, v := range want {
			vv := contents[k]
			if !slices.Equal(v, vv) {
				ztesting.AssertEqualSlice(t, "content not match for key="+k, v, contents[k])
			}
		}
		ztesting.AssertEqual(t, "returned data size not match", 2, len(contents))
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
}

func TestListFiles(t *testing.T) {
	t.Parallel()
	t.Run("contain empty", func(t *testing.T) {
		paths, err := zos.ListFiles(false, "")
		ztesting.AssertEqualSlice(t, "paths not match", []string{}, paths)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("invalid file path", func(t *testing.T) {
		paths, err := zos.ListFiles(false, "./not-found/")
		ztesting.AssertEqualSlice(t, "paths not match", []string{}, paths)
		ztesting.AssertEqual(t, "expected error not occurred", true, err != nil)
	})
	t.Run("not recursive", func(t *testing.T) {
		paths, err := zos.ListFiles(false, "./testdata/regular-file.txt", "./testdata/files/")
		want := []string{filepath.Clean("./testdata/regular-file.txt")}
		ztesting.AssertEqualSlice(t, "paths not match", want, paths)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
	t.Run("recursive", func(t *testing.T) {
		paths, err := zos.ListFiles(true, "./testdata/files/")
		want := []string{}
		want = append(want, filepath.Clean("./testdata/files/level0.txt"))
		want = append(want, filepath.Clean("./testdata/files/level1/level1.txt"))
		ztesting.AssertEqualSlice(t, "paths not match", want, paths)
		ztesting.AssertEqual(t, "non nil error returned", nil, err)
	})
}
