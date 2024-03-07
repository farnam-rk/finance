package dailyrotate

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// File log file daily structure
type File struct {
	sync.Mutex

	pathGenerator func(time.Time) string
	pathFormat    string

	Location *time.Location

	day     int
	path    string
	file    *os.File
	onClose func(path string, didRotate bool)

	lastWritePos int64
}

// close closes log file
func (f *File) close(didRotate bool) error {
	if f.file == nil {
		return nil
	}
	err := f.file.Close()
	f.file = nil
	if err == nil && f.onClose != nil {
		f.onClose(f.path, didRotate)
	}
	f.day = 0
	return err
}

// open opens log file
func (f *File) open() error {
	t := time.Now().UTC()
	if f.pathGenerator != nil {
		f.path = f.pathGenerator(t)
	} else {
		f.path = t.Format(f.pathFormat)
	}
	f.day = t.YearDay()

	dir := filepath.Dir(f.path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	flag := os.O_CREATE | os.O_WRONLY
	f.file, err = os.OpenFile(f.path, flag, 0644)
	if err != nil {
		return err
	}
	_, err = f.file.Seek(0, io.SeekEnd)
	return err
}

// reopenIfNeeded reopen log file in case of time change
func (f *File) reopenIfNeeded() error {
	t := time.Now().UTC()
	if t.YearDay() == f.day {
		return nil
	}
	err := f.close(true)
	if err != nil {
		return err
	}
	return f.open()
}

// NewFile generates new log file
func NewFile(pathFormat string, onClose func(path string, didRotate bool)) (*File, error) {
	return newFile(pathFormat, nil, onClose)
}

// NewFileWithPathGenerator generates new log file path
func NewFileWithPathGenerator(pathGenerator func(time.Time) string, onClose func(path string, didRotate bool)) (*File, error) {
	return newFile("", pathGenerator, onClose)
}

// newFile create new log file
func newFile(pathFormat string, pathGenerator func(time.Time) string, onClose func(path string, didRotate bool)) (*File, error) {
	f := &File{
		pathFormat:    pathFormat,
		pathGenerator: pathGenerator,
		Location:      time.UTC,
	}
	err := f.reopenIfNeeded()
	if err != nil {
		return nil, err
	}
	err = f.close(false)
	if err != nil {
		return nil, err
	}
	f.onClose = onClose
	return f, nil
}

// Close closes log file
func (f *File) Close() error {
	f.Lock()
	defer f.Unlock()
	return f.close(false)
}

// write writes to log file
func (f *File) write(d []byte, flush bool) (int64, int, error) {
	err := f.reopenIfNeeded()
	if err != nil {
		return 0, 0, err
	}
	f.lastWritePos, err = f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}
	n, err := f.file.Write(d)
	if err != nil {
		return 0, n, err
	}
	if flush {
		err = f.file.Sync()
	}
	return f.lastWritePos, n, err
}

// Write writes to log file
func (f *File) Write(d []byte) (int, error) {
	f.Lock()
	defer f.Unlock()
	_, n, err := f.write(d, false)
	return n, err
}

// Write2 writes to log file
func (f *File) Write2(d []byte, flush bool) (string, int64, int, error) {
	f.Lock()
	defer f.Unlock()
	writtenAtPos, n, err := f.write(d, flush)
	return f.path, writtenAtPos, n, err
}

// Flush flush file
func (f *File) Flush() error {
	f.Lock()
	defer f.Unlock()
	return f.file.Sync()
}
