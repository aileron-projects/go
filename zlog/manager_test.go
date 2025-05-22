package zlog

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/aileron-projects/go/ztesting"
)

func TestNewFileManager(t *testing.T) {
	t.Parallel()
	t.Run("pattern contains '/'", func(t *testing.T) {
		_, err := NewFileManager(&FileManagerConfig{
			Pattern: "foo/bar.txt",
		})
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("pattern contains '\\'", func(t *testing.T) {
		_, err := NewFileManager(&FileManagerConfig{
			Pattern: "foo\\bar.txt",
		})
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("invalid pattern", func(t *testing.T) {
		_, err := NewFileManager(&FileManagerConfig{
			Pattern: "bar.%x.%y.%z.txt",
		})
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("%i added", func(t *testing.T) {
		dir := t.TempDir()
		c := &FileManagerConfig{SrcDir: dir, DstDir: dir, Pattern: "test.txt"} // ".%i" should be added.
		_, err := NewFileManager(c)
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		ztesting.AssertEqual(t, "pattern not match", "test.txt.%i", c.Pattern)
	})
	t.Run("pattern contains specifier", func(t *testing.T) {
		dir := t.TempDir()
		specs := []string{"%H", "%U", "%G", "%p", "%P"}
		for _, s := range specs {
			c := &FileManagerConfig{SrcDir: dir, DstDir: dir, Pattern: s + ".%i"}
			_, err := NewFileManager(c)
			ztesting.AssertEqual(t, "error should be nil", nil, err)
			ztesting.AssertEqual(t, "pattern should have length", true, len(c.Pattern) > 0)
			ztesting.AssertEqual(t, "invalid pattern replace", true, c.Pattern[:len(c.Pattern)-3] != s)
		}
	})
	t.Run("srcDir create failed", func(t *testing.T) {
		dir := t.TempDir()
		c := &FileManagerConfig{SrcDir: dir + "/ng-\x00", DstDir: dir + "/ok"}
		_, err := NewFileManager(c)
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("dstDir create failed", func(t *testing.T) {
		dir := t.TempDir()
		c := &FileManagerConfig{SrcDir: dir + "/ok", DstDir: dir + "/ng-\x00"}
		_, err := NewFileManager(c)
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
}

func TestFileManager_NewFile(t *testing.T) {
	t.Parallel()
	m := &FileManager{srcDir: "testdata/", pattern: "test.%i.txt"}
	ztesting.AssertEqual(t, "file path not match", filepath.Join("testdata", "test.1.txt"), m.NewFile())
	ztesting.AssertEqual(t, "file path not match", filepath.Join("testdata", "test.2.txt"), m.NewFile())
	ztesting.AssertEqual(t, "file path not match", filepath.Join("testdata", "test.3.txt"), m.NewFile())
}

func TestFileManager_Manager(t *testing.T) {
	t.Parallel()
	t.Run("failed to list srcDir", func(t *testing.T) {
		m := &FileManager{srcDir: "not-found"}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("failed to archive", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/test.txt"
		os.WriteFile(path, []byte("testdata"), os.ModePerm)
		m := &FileManager{srcDir: dir, dstDir: "not-found", pattern: "test.txt"}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("failed to list dstDir", func(t *testing.T) {
		dir := t.TempDir()
		m := &FileManager{srcDir: dir, dstDir: "not-found", pattern: "test.txt"}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
	})
	t.Run("limit max history", func(t *testing.T) {
		dir := t.TempDir()
		copyDir("./testdata/manage/max-history", dir)
		m := &FileManager{srcDir: dir, dstDir: dir, pattern: "test.%i.txt", maxHistory: 3}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		want := []string{"test.5.txt", "test.6.txt", "test.9.txt"}
		ztesting.AssertEqual(t, "files not match", want, listDirFiles(dir))
	})
	t.Run("limit max history gzip", func(t *testing.T) {
		dir := t.TempDir()
		copyDir("./testdata/manage/max-history-gz", dir)
		m := &FileManager{srcDir: dir, dstDir: dir, pattern: "test.%i.txt", maxHistory: 3, gzipLv: gzip.BestSpeed}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		want := []string{"test.5.txt.gz", "test.6.txt.gz", "test.9.txt.gz"}
		ztesting.AssertEqual(t, "files not match", want, listDirFiles(dir))
	})
	t.Run("limit max age", func(t *testing.T) {
		dir := t.TempDir()
		now := time.Now().Unix()
		time1, time2 := strconv.FormatInt(now-100, 10), strconv.FormatInt(now-200, 10)
		time3, time4 := strconv.FormatInt(now-300, 10), strconv.FormatInt(now-400, 10)
		os.WriteFile(dir+"/test."+time1+".txt", []byte("testdata"), os.ModePerm)
		os.WriteFile(dir+"/test."+time2+".txt", []byte("testdata"), os.ModePerm)
		os.WriteFile(dir+"/test."+time3+".txt", []byte("testdata"), os.ModePerm)
		os.WriteFile(dir+"/test."+time4+".txt", []byte("testdata"), os.ModePerm)
		m := &FileManager{srcDir: dir, dstDir: dir, pattern: "test.%u.txt", maxAge: 250 * time.Second}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		want := []string{"test." + time2 + ".txt", "test." + time1 + ".txt"}
		ztesting.AssertEqual(t, "files not match", want, listDirFiles(dir))
	})
	t.Run("limit total size", func(t *testing.T) {
		dir := t.TempDir()
		copyDir("./testdata/manage/total-size", dir)
		m := &FileManager{srcDir: dir, dstDir: dir, pattern: "test.%i.txt", maxTotal: 25}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		want := []string{"test.4.txt", "test.5.txt"}
		ztesting.AssertEqual(t, "files not match", want, listDirFiles(dir))
	})
	t.Run("limit total size", func(t *testing.T) {
		dir := t.TempDir()
		copyDir("./testdata/manage/total-size-gz", dir)
		m := &FileManager{srcDir: dir, dstDir: dir, pattern: "test.%i.txt", maxTotal: 25, gzipLv: gzip.BestSpeed}
		err := m.Manage()
		ztesting.AssertEqual(t, "error should be nil", nil, err)
		want := []string{"test.4.txt.gz", "test.5.txt.gz"}
		ztesting.AssertEqual(t, "files not match", want, listDirFiles(dir))
	})
}

func listDirFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return err
		}
		srcF, _ := os.Open(path)
		defer srcF.Close()
		dstF, _ := os.Create(filepath.Join(dst, filepath.Base(path)))
		defer dstF.Close()
		_, err = io.Copy(dstF, srcF)
		return err
	})
}

func TestFileManager_archive(t *testing.T) {
	t.Parallel()
	t.Run("same dir", func(t *testing.T) {
		t.Run("with no compression", func(t *testing.T) {
			dir := t.TempDir()
			path := dir + "/test.txt"
			os.WriteFile(path, []byte("testdata"), os.ModePerm)
			m := &FileManager{srcDir: dir, dstDir: dir, gzipLv: gzip.NoCompression}
			err := m.archive([]*fileInfo{{name: "test.txt", path: path}})
			ztesting.AssertEqual(t, "error not match", nil, err)
			info, _ := os.Stat(path)
			ztesting.AssertEqual(t, "file type not match", true, info.Mode().IsRegular())
		})
		t.Run("with compression", func(t *testing.T) {
			dir := t.TempDir()
			path := dir + "/test.txt"
			os.WriteFile(path, []byte("testdata"), os.ModePerm)
			m := &FileManager{srcDir: dir, dstDir: dir, gzipLv: gzip.BestSpeed}
			err := m.archive([]*fileInfo{{name: "test.txt", path: path}})
			ztesting.AssertEqual(t, "error not match", nil, err)
			info, _ := os.Stat(path + ".gz")
			ztesting.AssertEqual(t, "file type not match", true, info.Mode().IsRegular())
		})
		t.Run("compression failed", func(t *testing.T) {
			dir := t.TempDir()
			path := dir + "/test.txt" // Not exist.
			m := &FileManager{srcDir: dir, dstDir: dir, gzipLv: gzip.BestSpeed}
			err := m.archive([]*fileInfo{{name: "test.txt", path: path}})
			ztesting.AssertEqual(t, "error should not be nil", true, err != nil)
		})
	})
	t.Run("different dir", func(t *testing.T) {
		t.Run("with no compression", func(t *testing.T) {
			srcDir := t.TempDir()
			dstDir := t.TempDir()
			path := srcDir + "/test.txt"
			os.WriteFile(path, []byte("testdata"), os.ModePerm)
			m := &FileManager{srcDir: srcDir, dstDir: dstDir, gzipLv: gzip.NoCompression}
			err := m.archive([]*fileInfo{{name: "test.txt", path: path}})
			ztesting.AssertEqual(t, "error not match", nil, err)
			info, _ := os.Stat(dstDir + "/test.txt")
			ztesting.AssertEqual(t, "file type not match", true, info.Mode().IsRegular())
		})
		t.Run("with compression", func(t *testing.T) {
			srcDir := t.TempDir()
			dstDir := t.TempDir()
			path := srcDir + "/test.txt"
			os.WriteFile(path, []byte("testdata"), os.ModePerm)
			m := &FileManager{srcDir: srcDir, dstDir: dstDir, gzipLv: gzip.BestSpeed}
			err := m.archive([]*fileInfo{{name: "test.txt", path: path}})
			ztesting.AssertEqual(t, "error not match", nil, err)
			info, _ := os.Stat(dstDir + "/test.txt.gz")
			ztesting.AssertEqual(t, "file type not match", true, info.Mode().IsRegular())
		})
	})
}

func TestGzipCompressFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir() // Temp directory for this test.
	srcFile := tmp + "/test.txt"
	dstFile := tmp + "/compressed.gz"

	testCases := map[string]struct {
		src, dst string
		level    int
		err      error
	}{
		"invalid src":      {src: "not found", err: errors.New("placeholder")},
		"invalid dst":      {src: "testdata/gzip.txt", dst: "invalid-\x00", err: errors.New("placeholder")},
		"no compression":   {src: srcFile, dst: dstFile, level: gzip.NoCompression},
		"huffman only":     {src: srcFile, dst: dstFile, level: gzip.HuffmanOnly},
		"best compression": {src: srcFile, dst: dstFile, level: gzip.BestCompression},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			os.Remove(dstFile)
			os.WriteFile(srcFile, []byte("testdata"), os.ModePerm)
			err := gzipCompressFile(tc.src, tc.dst, tc.level)
			if tc.err != nil {
				// Currently we only check if the error is nil or not because the PathError
				// contains platform dependent errors.
				ztesting.AssertEqual(t, "error not match", true, err != nil)
				return
			}
			ztesting.AssertEqual(t, "error not match", nil, err)
			b, err := os.ReadFile(dstFile)
			r, _ := gzip.NewReader(bytes.NewReader(b))
			bb, _ := io.ReadAll(r)
			ztesting.AssertEqual(t, "error not match", nil, err)
			ztesting.AssertEqual(t, "content not match", "testdata", string(bb))
		})
	}
}

func TestListFiles(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		pattern, dir string
		files        []*fileInfo
		err          error
	}{
		"invalid dir": {dir: "not-found", err: errors.New("placeholder")},
		"sort by index": {
			dir:     "testdata/listfiles/",
			pattern: "sort-by-index.%i.txt",
			files: []*fileInfo{
				{name: "sort-by-index.03.txt", index: 3},
				{name: "sort-by-index.02.txt", index: 2},
				{name: "sort-by-index.01.txt", index: 1},
			},
		},
		"sort by timestamp": {
			dir:     "testdata/listfiles/",
			pattern: "sort-by-time.%Y-%M-%D.txt",
			files: []*fileInfo{
				{name: "sort-by-time.2025-02-28.txt", created: time.Date(2025, 02, 28, 0, 0, 0, 0, time.Local)},
				{name: "sort-by-time.2025-02-01.txt", created: time.Date(2025, 02, 01, 0, 0, 0, 0, time.Local)},
				{name: "sort-by-time.2025-01-31.txt", created: time.Date(2025, 01, 31, 0, 0, 0, 0, time.Local)},
				{name: "sort-by-time.2025-01-01.txt", created: time.Date(2025, 01, 01, 0, 0, 0, 0, time.Local)},
			},
		},
		"sort by unix": {
			dir:     "testdata/listfiles/",
			pattern: "sort-by-unix.%u.txt",
			files: []*fileInfo{
				{name: "sort-by-unix.1800000100.txt", created: time.Unix(1800000100, 0)},
				{name: "sort-by-unix.1800000010.txt", created: time.Unix(1800000010, 0)},
				{name: "sort-by-unix.1800000000.txt", created: time.Unix(1800000000, 0)},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			files, err := listFiles(tc.pattern, tc.dir)
			if tc.err != nil {
				// Currently we only check if the error is nil or not because the PathError
				// contains platform dependent errors.
				ztesting.AssertEqual(t, "error not match", true, err != nil)
				return
			}
			ztesting.AssertEqual(t, "error not match", nil, err)
			ztesting.AssertEqual(t, "file length not match", len(tc.files), len(files))
			for i := range tc.files {
				ztesting.AssertEqual(t, "file name not match", tc.files[i].name, files[i].name)
				ztesting.AssertEqual(t, "index not match", tc.files[i].index, files[i].index)
				ztesting.AssertEqual(t, "created time not match", tc.files[i].created, files[i].created)
			}
		})
	}
}

func TestFormatFileName(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		format, str string
		t           time.Time
		index       int
		ok          bool
	}{
		"Y only":        {"%Y", "2025", time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local), 0, true},
		"M only":        {"%M", "12", time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local), 0, true},
		"D only":        {"%D", "31", time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local), 0, true},
		"h only":        {"%h", "23", time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local), 0, true},
		"m only":        {"%m", "59", time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local), 0, true},
		"s only":        {"%m", "59", time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local), 0, true},
		"u only":        {"%u", "1234567890", time.Unix(1234567890, 0), 0, true},
		"i only":        {"%i", "99", time.Time{}, 99, true},
		"YMD":           {"%Y%M%D", "20250423", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"YMD-i":         {"%Y%M%D-%i", "20250423-123", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 123, true},
		"hms":           {"%h%m%m", "235959", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"hms-i":         {"%h%m%m-%i", "235959-123", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 123, true},
		"Y-M-D":         {"%Y-%M-%D", "2025-04-23", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"Y_M_D":         {"%Y_%M_%D", "2025_04_23", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"YMD-hms":       {"%Y%M%D-%h%m%s", "20250423-235959", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"YMD_hms":       {"%Y%M%D_%h%m%s", "20250423_235959", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"contain Y-M-D": {"test-%Y-%M-%D.log", "test-2025-04-23.log", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"contain i":     {"test-%i.log", "test-1234.log", time.Time{}, 1234, true},
		"no spec":       {"test", "test", time.Time{}, 0, true},
		"undefined":     {"%x", "", time.Time{}, 0, false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			str, ok := formatFileName(tc.format, tc.t, tc.index)
			ztesting.AssertEqual(t, "string not match", tc.str, str)
			ztesting.AssertEqual(t, "ok not match", tc.ok, ok)
		})
	}
}

func TestScanFileName(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		format, str string
		t           time.Time
		index       int
		ok          bool
	}{
		"Y only":        {"%Y", "2025", time.Time{}, 0, true},
		"M only":        {"%M", "12", time.Time{}, 0, true},
		"D only":        {"%D", "31", time.Time{}, 0, true},
		"h only":        {"%h", "23", time.Time{}, 0, true},
		"m only":        {"%m", "59", time.Time{}, 0, true},
		"s only":        {"%m", "59", time.Time{}, 0, true},
		"u only":        {"%u", "1234567890", time.Unix(1234567890, 0), 0, true},
		"i only":        {"%i", "99", time.Time{}, 99, true},
		"Y invalid":     {"%Y", "abcd", time.Time{}, 0, false},
		"M invalid":     {"%M", "13", time.Time{}, 0, false},
		"D invalid":     {"%D", "32", time.Time{}, 0, false},
		"h invalid":     {"%h", "24", time.Time{}, 0, false},
		"m invalid":     {"%m", "60", time.Time{}, 0, false},
		"s invalid":     {"%m", "60", time.Time{}, 0, false},
		"u invalid":     {"%u", "xx", time.Time{}, 0, false},
		"i invalid":     {"%i", "xx", time.Time{}, 0, false},
		"YMD":           {"%Y%M%D", "20250423", time.Date(2025, 4, 23, 0, 0, 0, 0, time.Local), 0, true},
		"YMD-i":         {"%Y%M%D-%i", "20250423-123", time.Date(2025, 4, 23, 0, 0, 0, 0, time.Local), 123, true},
		"hms":           {"%h%m%m", "235959", time.Time{}, 0, true},
		"hms-i":         {"%h%m%m-%i", "235959-123", time.Time{}, 123, true},
		"Y-M-D":         {"%Y-%M-%D", "2025-04-23", time.Date(2025, 4, 23, 0, 0, 0, 0, time.Local), 0, true},
		"Y_M_D":         {"%Y_%M_%D", "2025_04_23", time.Date(2025, 4, 23, 0, 0, 0, 0, time.Local), 0, true},
		"YMD-hms":       {"%Y%M%D-%h%m%s", "20250423-235959", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"YMD_hms":       {"%Y%M%D_%h%m%s", "20250423_235959", time.Date(2025, 4, 23, 23, 59, 59, 0, time.Local), 0, true},
		"contain Y-M-D": {"test-%Y-%M-%D.log", "test-2025-04-23.log", time.Date(2025, 4, 23, 0, 0, 0, 0, time.Local), 0, true},
		"contain i":     {"test-%i.log", "test-1234.log", time.Time{}, 1234, true},
		"no spec 1":     {"test", "test", time.Time{}, 0, true},
		"no spec 2":     {"test", "hello", time.Time{}, 0, false},
		"not match":     {"test%i.log", "hello123.log", time.Time{}, 0, false},
		"undefined":     {"%x", "xx", time.Time{}, 0, false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tm, index, ok := scanFileName(tc.format, tc.str)
			ztesting.AssertEqual(t, "time not match", tc.t, tm)
			ztesting.AssertEqual(t, "index not match", tc.index, index)
			ztesting.AssertEqual(t, "char not match", tc.ok, ok)
		})
	}
}

func TestScanFormat(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		format       string
		prefix, rest string
		char         byte
	}{
		"case01": {"", "", "", 0x00},
		"case02": {"test", "test", "", 0x00},
		"case03": {"%test", "", "est", 't'},
		"case04": {"test%", "test%", "", 0x00},
		"case05": {"te%st", "te", "t", 's'},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			prefix, rest, char := scanFormat(tc.format)
			ztesting.AssertEqual(t, "prefix not match", tc.prefix, prefix)
			ztesting.AssertEqual(t, "rest string not match", tc.rest, rest)
			ztesting.AssertEqual(t, "char not match", tc.char, char)
		})
	}
}

func TestScanNumber(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		str   string
		digit int
		num   int64
		rest  string
		ok    bool
	}{
		"case01": {"", 0, 0, "", false},
		"case02": {"0123", 0, 0, "0123", false},
		"case03": {"0123", 1, 0, "123", true},
		"case04": {"0123", 2, 1, "23", true},
		"case05": {"0123", 3, 12, "3", true},
		"case06": {"0123", 4, 123, "", true},
		"case07": {"0123", 5, 0, "", false}, // Digits over.
		"case08": {"0123", -1, 123, "", true},
		"case09": {"test", -1, 0, "test", false},
		"case10": {"test", 1, 0, "est", false},
		"case11": {"0123test", 3, 12, "3test", true},
		"case12": {"0123test", 5, 0, "est", false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			num, rest, ok := scanNumber(tc.str, tc.digit)
			ztesting.AssertEqual(t, "parsed num not match", tc.num, num)
			ztesting.AssertEqual(t, "rest string not match", tc.rest, rest)
			ztesting.AssertEqual(t, "ok not match", tc.ok, ok)
		})
	}
}
