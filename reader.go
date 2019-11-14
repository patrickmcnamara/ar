package ar

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Reader provides sequential access to an ar archive. Reader.NextFile advances
// to the next file in the archive, including the first. The Reader can then be
// be used as an io.Reader to access the file's data.
type Reader struct {
	r  io.Reader
	nb int
	nl bool
}

// NewReader creates a new Reader reading from r.
//
// It will return io.ErrUnexpectedEOF or ErrMissingMagic if the first eight
// bytes do not match ARMAG.
func NewReader(r io.Reader) (*Reader, error) {
	sig := make([]byte, 8)
	if _, err := r.Read(sig); err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return nil, err
	}
	if !bytes.Equal(sig, []byte(ARMAG)) {
		return nil, ErrMissingMagic
	}
	return &Reader{r: r}, nil
}

// NextFile advances to the next file in the ar archive, including the first.
// The hdr.Size then determines how much can be read from the Reader to access
// the file's content.
//
// If the current file has not been fully read from the Reader, the rest of the
// file is discarded.
//
// io.EOF is returned if there are no more files.
func (r *Reader) NextFile() (hdr Header, err error) {
	// discard previous file
	if r.nb > 0 {
		if _, err = io.Copy(ioutil.Discard, r); err != nil {
			return
		}
	}

	// discard possible newline
	if r.nl {
		if _, err = io.CopyN(ioutil.Discard, r.r, 1); err != nil {
			return
		}
	}

	// read whole header
	buf := make([]byte, 60)
	var n int
	if n, err = r.r.Read(buf); err == io.EOF {
		if n != 0 {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	if err != nil {
		return
	}

	// name
	hdr.Name = strings.TrimSuffix(strings.TrimRight(string(buf[:16]), " "), "/")

	// date
	date, _ := strconv.Atoi(string(bytes.TrimRight(buf[16:28], " ")))
	hdr.Date = time.Unix(int64(date), 0)

	// UID
	uid, _ := strconv.ParseUint(string(bytes.TrimRight(buf[28:34], " ")), 10, 32)
	hdr.UID = uint32(uid)

	// GID
	gid, _ := strconv.ParseUint(string(bytes.TrimRight(buf[34:42], " ")), 10, 32)
	hdr.GID = uint32(gid)

	// mode
	mode, _ := strconv.ParseUint(string(bytes.TrimRight(buf[40:48], " ")), 8, 32)
	hdr.Mode = os.FileMode(mode)

	// size
	size, _ := strconv.Atoi(string(bytes.TrimRight(buf[48:58], " ")))
	hdr.Size = int(size)

	// ignore last two bytes, ARFMAG
	_ = buf[58:]

	// set up writer and next file
	r.nb = hdr.Size
	r.nl = hdr.Size&1 != 0

	return
}

// Read reads from the current file in the ar archive. It will return io.EOF
// when it reaches the end of the current file.
func (r *Reader) Read(p []byte) (n int, err error) {
	// already read all of file
	if r.nb == 0 {
		err = io.EOF
		return
	}

	// read upto nb
	nb := len(p)
	if nb > r.nb {
		nb = r.nb
	}

	// read from underlying reader
	n, err = r.r.Read(p[:nb])
	if err != nil {
		return
	}

	// decrease read from remaining
	r.nb -= n

	return
}
