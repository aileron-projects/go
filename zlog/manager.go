package zlog

import (
	"bytes"
	"cmp"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FileManagerConfig is the configuration for the [FileManager].
// Use [NewFileManager] to create a new instance of it.
type FileManagerConfig struct {
	// MaxAge is the maximum age to keep the archived files.
	// It does not work when time specifiers are not exit
	// in the Pattern. At least %Y, %M and %D, or %u must be specified.
	// Zero or negative value means no limit.
	MaxAge time.Duration
	// MaxHistory is the maximum number of archived files to keep.
	// Zero or negative value means no limit.
	MaxHistory int
	// MaxTotalBytes is the maximum total byte size to keep archived files.
	// Zero or negative value means no limit.
	MaxTotalBytes int64
	// GzipLv specifies the gzip compression level.
	// If non-zero, gzip compression is applied to archived files,
	// and the ".gz" file extension is added.
	// Valid range is from [compress/gzip.HuffmanOnly] to [compress/gzip.BestCompression].
	// If zero, compression is disabled.
	GzipLv int
	// SrcDir is the directory path where source files are exists.
	// File names should have the pattern specified at the Pattern.
	// Current working directory is used when empty.
	SrcDir string
	// DesDir is the destination directory path
	// to place archived files.
	// File names will have the pattern specified at the Pattern
	// with or without additional ".gz" extension.
	// Current working directory is used when empty.
	DstDir string
	// Pattern is the file name pattern to be manged.
	// If empty, "application.%i.log" is used.
	// Following specifiers can be used in the pattern.
	// A ".%i" will be added when no specifiers are found.
	//	%Y : YYYY 4 digits year. 0 <= YYYY
	//	%M : MM 2 digits month. 1 <= MM <= 12
	//	%D : DD 2 digits day of month. 1 <= DD <= 31
	//	%h : hh 2 digits hour. 0 <= hh <= 23
	//	%m : mm 2 digits minute. 0 <= mm <= 59
	//	%s : ss 2 digits second. 0 <= ss <= 59
	//	%u : unix second with free digits. 0 <= unix
	//	%i : index with free digits. 0 <= index
	//	%H : hostname
	//	%U : user id. "-1" on windows.
	//	%G : user group id. "-1" on windows.
	//	%p : pid (process id)
	//	%P : ppid (parent process id)
	Pattern string
}

// NewFileManager returns a new instance of [FileManager].
func NewFileManager(c *FileManagerConfig) (*FileManager, error) {
	c = cmp.Or(c, &FileManagerConfig{})
	c.GzipLv = max(gzip.HuffmanOnly, min(c.GzipLv, gzip.BestCompression))
	if strings.Contains(c.Pattern, "/") || strings.Contains(c.Pattern, "\\") {
		return nil, errors.New("zlog: pattern should not contain directory path")
	}
	// Replace specifiers with fixed value.
	hostname, _ := os.Hostname()
	c.Pattern = strings.ReplaceAll(c.Pattern, "%H", hostname)
	c.Pattern = strings.ReplaceAll(c.Pattern, "%U", strconv.Itoa(os.Getuid()))  // -1 on windows.
	c.Pattern = strings.ReplaceAll(c.Pattern, "%G", strconv.Itoa(os.Getegid())) // -1 on windows.
	c.Pattern = strings.ReplaceAll(c.Pattern, "%p", strconv.Itoa(os.Getpid()))
	c.Pattern = strings.ReplaceAll(c.Pattern, "%P", strconv.Itoa(os.Getppid()))
	c.Pattern = cmp.Or(c.Pattern, "application.%i.log") // Default if empty.
	if _, ok := formatFileName(c.Pattern, time.Time{}, 123); !ok {
		return nil, errors.New("zlog: invalid specifier in pattern")
	}
	_, _, char := scanFormat(c.Pattern)
	if char == 0x00 || !bytes.ContainsAny([]byte{char}, "YMDhmsui") {
		c.Pattern += ".%i"
	}
	c.SrcDir = filepath.Clean(c.SrcDir)
	c.DstDir = filepath.Clean(c.DstDir)
	if err := os.MkdirAll(c.SrcDir, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(c.DstDir, os.ModePerm); err != nil {
		return nil, err
	}
	m := &FileManager{
		maxAge:     c.MaxAge,
		maxHistory: c.MaxHistory,
		maxTotal:   c.MaxTotalBytes,
		gzipLv:     c.GzipLv,
		srcDir:     c.SrcDir,
		dstDir:     c.DstDir,
		pattern:    c.Pattern,
	}
	return m, m.Manage()
}

// FileManager manages files that have the specified format
// in a directory. It archives file from srcDir to dstDir and
// manages their life.
// Use [NewFileManager] to instantiate [FileManager].
//
// FileManager has the following features:
//
//   - maxAge: Remove archived files older than the age.
//   - maxHistory: Limit the number of archived files.
//   - maxTotalSize: Limit the total size of archived files.
//   - gzip: Compress archived files.
type FileManager struct {
	mu sync.Mutex
	// maxAge is the maximum age of the backup files in second.
	// Files older than this age will be removed.
	// Zero or negative means no limitation.
	maxAge time.Duration
	// maxBackup is the maximum number of the backup files.
	// If the number of backup files exceeded this value,
	// backup files will be removed from the older one.
	// Zero or negative means no limitation.
	maxHistory int
	// maxTotal is the maximum total file size in bytes.
	// Zero or negative means no limitation.
	maxTotal int64
	// gzipLv is the gzip compression level.
	// If not zero, files are gzip compressed when archiving.
	gzipLv int
	// srcDir is the source directory.
	// Files will be moved from srcDir to dstDir when archiving them.
	// Current working directory is used if empty.
	srcDir string
	// dstDir is the destination directory.
	// Files will be moved from srcDir to dstDir when archiving them.
	// Current working directory is used if empty.
	dstDir string
	// pattern is the file name pattern.
	// pattern may not contain '/'.
	pattern string
	// index is the current index number.
	// This won't be reset otherwise the instance of
	// FileManage had been recreated.
	index int
}

// NewFile returns a new file path that can be managed by the FileManager.
// The FileManager does not create the file returned by the NewFile.
// Callers should open, create, close or remove the file themselves.
// Note that calling NewFile increments the internal index that may be
// used in the file name pattern for the value of '%i'.
func (m *FileManager) NewFile() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.index++
	index := m.index
	name, _ := formatFileName(m.pattern, time.Now(), index) // pattern is already validated.
	return filepath.Join(m.srcDir, name)
}

// Manage manages archived files.
// Manage may take longer time depending on the number of files
// or size of files to be compressed. Manage is safe for concurrent call.
func (m *FileManager) Manage() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	files, err := listFiles(m.pattern, m.srcDir)
	if err != nil {
		return err
	}
	if err := m.archive(files); err != nil {
		return err
	}
	ext := ""
	if m.gzipLv != gzip.NoCompression {
		ext = ".gz"
	}
	files, err = listFiles(m.pattern+ext, m.dstDir)
	if err != nil {
		return err
	}
	if len(files) > 0 && files[0].index >= m.index {
		m.index = files[0].index // Update current index.
	}
	totalSize := int64(0)
	var errs []error
	for i, file := range files {
		if m.maxHistory > 0 && i >= m.maxHistory {
			errs = appendNonNil(errs, os.Remove(file.path))
		}
		if m.maxAge > 0 && file.age > m.maxAge {
			errs = appendNonNil(errs, os.Remove(file.path))
		}
		totalSize += file.size
		if m.maxTotal > 0 && totalSize > m.maxTotal {
			errs = appendNonNil(errs, os.Remove(file.path))
		}
	}
	return errors.Join(errs...)
}

// archive archives file from srcDir to dstDir.
// Given files must be present in the srcDir.
func (m *FileManager) archive(files []*fileInfo) error {
	if m.srcDir == m.dstDir && m.gzipLv == gzip.NoCompression {
		// No need to rename.
		// No need to compress.
		return nil
	}
	var errs []error
	for _, file := range files {
		if m.gzipLv == gzip.NoCompression {
			errs = appendNonNil(errs, os.Rename(file.path, filepath.Join(m.dstDir, file.name)))
			continue
		}
		src := file.path
		dst := filepath.Join(m.dstDir, file.name+".gz") // On windows, '\' is used.
		errs = appendNonNil(errs, gzipCompressFile(src, dst, m.gzipLv))
	}
	return errors.Join(errs...)
}

func appendNonNil(errs []error, err error) []error {
	if err == nil {
		return errs
	}
	return append(errs, err)
}

// fileInfo is the file information of archived files
// that are managed by the [FileManager].
type fileInfo struct {
	name    string        // Name is the file name. '/' not contained.
	path    string        // Path is the file path. '/' may be contained.
	size    int64         // Size is the file size in bytes.
	created time.Time     // Created is the created timestamp if any.
	age     time.Duration // Age is the file age based on the create time if available.
	index   int           // Index is the incremental index number if any.
}

// gzipCompress compressed src files to the dst file.
// src file will be removed after compression successfully finished.
// Gzip compression level MUST be in appropriate range.
// See also [gzip] package.
func gzipCompressFile(src, dst string, level int) error {
	srcFile, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	level = max(gzip.HuffmanOnly, min(level, gzip.BestCompression))
	gw, _ := gzip.NewWriterLevel(dstFile, level)
	defer gw.Close()

	if _, err := io.CopyBuffer(gw, srcFile, make([]byte, 32<<10)); err != nil {
		return err
	}
	// Close and remove the source file.
	srcFile.Close()
	return os.Remove(src)
}

// listFiles returns file info in the dir those file names
// match to the given pattern.
// It does not search files of sub-directories of the dir.
// Returned slice of fileInfo is sorted.
func listFiles(pattern, dir string) ([]*fileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	now := time.Now().Local()
	files := make([]*fileInfo, 0, 8) // Use initial size 8.
	for _, entry := range entries {
		info, err := entry.Info()
		if !entry.Type().IsRegular() || err != nil {
			continue
		}
		created, index, ok := scanFileName(pattern, entry.Name())
		if !ok {
			continue
		}
		var age time.Duration
		if created.Year() > 0 && created.Month() > 0 && created.Day() > 0 {
			age = now.Sub(created) // Accept age only when %Y,%M,%D are valid.
		}
		files = append(files, &fileInfo{
			name:    entry.Name(),
			path:    filepath.Join(dir, entry.Name()), // On windows, '\' is used.
			size:    info.Size(),
			created: created,
			age:     age,
			index:   index,
		})
	}
	sort.SliceStable(files, func(i, j int) bool {
		if files[i].age != files[j].age {
			return files[i].age < files[j].age
		}
		return files[i].index > files[j].index
	})
	return files, nil
}

// formatFileName returns formatted string.
// Use [scanFileName] to parse time and index from file names.
// Timezone is always [time.Local].
//
// Allowed format specifiers are:
//
//	%Y : YYYY 4 digits year. 0 <= YYYY
//	%M : MM 2 digits month. 1 <= MM <= 12
//	%D : DD 2 digits day of month. 1 <= DD <= 31
//	%h : hh 2 digits hour. 0 <= hh <= 23
//	%m : mm 2 digits minute. 0 <= mm <= 59
//	%s : ss 2 digits second. 0 <= ss <= 59
//	%u : unix second with free digits. 0 <= unix
//	%i : index with free digits. 0 <= index
func formatFileName(format string, t time.Time, index int) (str string, ok bool) {
	var builder strings.Builder
	builder.Grow(len(format) + 5) // Add +5 just in case.
	var prefix string
	var char byte
	t = t.Local() // Force local time.
	year, month, day := t.Date()
	for format != "" {
		prefix, format, char = scanFormat(format)
		_, _ = builder.WriteString(prefix)
		if char == 0x00 {
			return builder.String(), format == ""
		}
		switch char {
		case 'Y':
			_, _ = builder.WriteString(fmt.Sprintf("%04d", year))
		case 'M':
			_, _ = builder.WriteString(fmt.Sprintf("%02d", month))
		case 'D':
			_, _ = builder.WriteString(fmt.Sprintf("%02d", day))
		case 'h':
			_, _ = builder.WriteString(fmt.Sprintf("%02d", t.Hour()))
		case 'm':
			_, _ = builder.WriteString(fmt.Sprintf("%02d", t.Minute()))
		case 's':
			_, _ = builder.WriteString(fmt.Sprintf("%02d", t.Second()))
		case 'u':
			_, _ = builder.WriteString(strconv.FormatInt(t.Unix(), 10))
		case 'i':
			_, _ = builder.WriteString(strconv.Itoa(index))
		default:
			return "", false // Invalid '%'.
		}
	}
	return builder.String(), true
}

// scanFileName scans str with the given format.
// It parses timestamp info and index from the str if format specifiers
// are exist in the format.
// It returns false when the str does not comply with the format.
// It returns parsed time and index with true when the str complies with the format.
// Use [formatFileName] to generate file names that comply with the format.
// Time zone of parsed time is always [time.Local].
// If both time '%u' and other time specifiers are exist in the format,
// it uses unix time as the returned time t.
// It identifies time that are not exist in the format as zero.
// For example, when the format has '%Y-%M-%D', the other time of
// '%h', '%m' and '%s' are recognized as zero.
// Returned time always be zero value if '%Y', '%M' or 'D' not exist.
// See also https://www.w3.org/TR/NOTE-datetime.
//
// Allowed format specifiers are:
//
//	%Y : YYYY 4 digits year. 0 <= YYYY
//	%M : MM 2 digits month. 1 <= MM <= 12
//	%D : DD 2 digits day of month. 1 <= DD <= 31
//	%h : hh 2 digits hour. 0 <= hh <= 23
//	%m : mm 2 digits minute. 0 <= mm <= 59
//	%s : ss 2 digits second. 0 <= ss <= 59
//	%u : unix second with free digits. 0 <= unix
//	%i : index with free digits. 0 <= index
func scanFileName(format, str string) (t time.Time, index int, ok bool) {
	var prefix string
	var char byte
	var year, month, day, hour, minute, second, unix, idx int64
	for format != "" {
		prefix, format, char = scanFormat(format)
		if !strings.HasPrefix(str, prefix) {
			return time.Time{}, 0, false
		}
		str = strings.TrimPrefix(str, prefix)
		if char == 0x00 {
			break
		}
		var ok bool
		switch char {
		case 'Y':
			year, str, ok = scanNumber(str, 4)
		case 'M':
			month, str, ok = scanNumber(str, 2)
		case 'D':
			day, str, ok = scanNumber(str, 2)
		case 'h':
			hour, str, ok = scanNumber(str, 2)
		case 'm':
			minute, str, ok = scanNumber(str, 2)
		case 's':
			second, str, ok = scanNumber(str, 2)
		case 'u':
			unix, str, ok = scanNumber(str, -1) // Free digits.
		case 'i':
			idx, str, ok = scanNumber(str, -1) // Free digits.
			index = int(idx)
		default:
			ok = false // Invalid '%'.
		}
		if !ok {
			return time.Time{}, 0, false
		}
	}
	if month > 12 || day > 31 || hour > 23 || minute > 59 || second > 59 {
		return time.Time{}, 0, false // Invalid number.
	}
	if unix > 0 {
		t = time.Unix(unix, 0)
	} else if year > 0 && month > 0 && day > 0 {
		t = time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.Local)
	}
	return t, index, str == "" && format == ""
}

// scanFormat scans the format string and identifies the first format specifier
// indicated by '%'. It returns the string before the specifier as 'prefix',
// the string after the specifier as 'rest', and the format specifier itself as 'char'.
//
// For example, given the string "foo%dbar", it returns:
//
//	prefix: "foo"
//	rest:   "bar"
//	char:   'd'
//
// If no format specifier is found, the entire format string is returned as 'prefix',
// with an empty 'rest' and 0x00 as 'char'.
// Note that a '%' at the end of the string is not treated as a format specifier.
func scanFormat(format string) (prefix, rest string, char byte) {
	n := len(format)
	for i := range format {
		if format[i] == '%' && i < n-1 {
			return format[:i], format[i+2:], format[i+1]
		}
	}
	return format, "", 0x00
}

// scanNumber scans numbers with the specified digit from the beginning of the str.
// It returns parsed number and the rest of str.
// It returns false when no numbers found at the prefix of the str or
// len(str) is less than digit.
// When digit<0, scanNumber scans number until it encounters non numeric character.
// digit=0 always results in false because an empty string cannot be converted into a number.
// It uses [strconv.Atoi] to parse number string to integer.
//
// For example:
//
//	scanNumber("012alice", 3) -> 12, "alice", true
//	scanNumber("012alice", 2) -> 1, "2alice", true
//	scanNumber("alice012", 2) -> 0, "", false
//	scanNumber("0123alice", -1) -> 123, "alice", true
func scanNumber(str string, digit int) (num int64, rest string, ok bool) {
	var numStr string
	if digit >= 0 {
		if len(str) < digit {
			return 0, "", false
		}
		numStr = str[:digit]
		rest = str[digit:]
	} else {
		numStr = str
		for i := range str {
			if '0' <= str[i] && str[i] <= '9' {
				continue
			}
			numStr = str[:i]
			rest = str[i:]
			break
		}
	}
	// Format "123" or "0123" are allowed.
	// Empty string results in an error.
	num, err := strconv.ParseInt(numStr, 10, 64)
	return num, rest, err == nil
}
