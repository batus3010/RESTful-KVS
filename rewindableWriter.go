package kvs

import "io"

// rewindableWriter encapsulate the action: "when we write we go from the beginning" of the file
type rewindableWriter struct {
	file io.ReadWriteSeeker
}

func (w *rewindableWriter) Write(p []byte) (n int, err error) {
	// rewind to the beginning
	if _, err = w.file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}
	// truncate if possible
	if t, ok := w.file.(interface{ Truncate(int64) error }); ok {
		if err = t.Truncate(0); err != nil {
			return 0, err
		}
		// make sure cursor is at zero after truncation
		if _, err = w.file.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
	}
	// now write the full payload in one shot
	return w.file.Write(p)
}
