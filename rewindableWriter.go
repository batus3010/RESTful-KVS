package kvs

import (
	"io"
	"os"
)

// rewindableWriter encapsulate the action: "when we write we go from the beginning" of the file
type rewindableWriter struct {
	file *os.File
}

func (w *rewindableWriter) Write(p []byte) (n int, err error) {
	// 1) Rewind
	if _, err = w.file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	// 2) Truncate the file to zero length
	if err = w.file.Truncate(0); err != nil {
		return 0, err
	}
	// 3) Rewind again (Truncate may move the offset)
	if _, err = w.file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	// 4) Write new data
	return w.file.Write(p)
}
