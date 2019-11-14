package ar

import (
	"errors"
	"io"
	"strconv"
)

// Writer provides sequential writing of an ar archive. Writer.NextFile begins a
// new file. The Writer can then be used as an io.Writer to supply that file's
// content.
type Writer struct {
	w  io.Writer
	nb int
	nl bool
}

// NewWriter a new Writer writing to w.
func NewWriter(w io.Writer) (*Writer, error) {
	// write initial file sig and return
	_, err := w.Write([]byte(ARMAG))
	return &Writer{w: w}, err
}

// NextFile writes h and prepares the Writer to accept the file's content. The
// hdr.Size determines how many bytes can be written to the Writer.
//
// It will return an ErrFileNotDone if the previous file was not fully written.
func (w *Writer) NextFile(h Header) error {
	// previous file not done
	if w.nb > 0 {
		return ErrFileNotDone
	}

	// header buffer
	buf := make([]byte, 60)

	// name
	if len(h.Name) > 16 {
		return errors.New("ar: name too long")
	}
	if h.Name == "" {
		return errors.New("ar: name cannot be empty")
	}
	copy(buf, h.Name)

	// date
	date := strconv.Itoa(int(h.Date.Unix()))
	if len(date) > 12 {
		return errors.New("ar: date too long")
	}
	copy(buf[16:], date)

	// UID
	uid := strconv.Itoa(int(h.UID))
	if len(uid) > 6 {
		return errors.New("ar: UID too long")
	}
	copy(buf[28:], uid)

	// GID
	gid := strconv.Itoa(int(h.GID))
	if len(gid) > 6 {
		return errors.New("ar: GID too long")
	}
	copy(buf[34:], gid)

	// mode
	mode := strconv.FormatUint(uint64(h.Mode), 8)
	if len(mode) > 8 {
		return errors.New("ar: mode too long")
	}
	copy(buf[40:], mode)

	// size
	size := strconv.Itoa(int(h.Size))
	if len(size) > 10 {
		return errors.New("ar: file too big")
	}
	copy(buf[48:], size)

	// fmag
	copy(buf[58:], ARFMAG)

	// replace zeroes with spaces
	for i, val := range buf {
		if val == 0 {
			buf[i] = ' '
		}
	}

	// set up writer
	w.nb = h.Size
	w.nl = h.Size&1 != 0

	// write to underlying writer and return
	_, err := w.w.Write(buf)
	return err
}

// Write writes to the current file in the ar archive.
//
// It will return ErrWriteTooLong if more than Header.Size writes are written
// since the previous Writer.NextFile.
func (w *Writer) Write(p []byte) (n int, err error) {
	// write upto nb
	nb := len(p)
	if nb > w.nb {
		nb = w.nb
	}

	// write to underlying writer
	n, err = w.w.Write(p[:nb])
	if err != nil {
		return
	}

	// decrease written from remaining
	w.nb -= n

	// append extra newline if required and return
	if w.nl && w.nb == 0 {
		_, err = w.w.Write([]byte{'\n'})
		if err != nil {
			return
		}
	}

	// write too long
	if nb != len(p) {
		err = ErrWriteTooLong
	}

	return
}
