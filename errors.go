package ar

import "errors"

// Reader errors.
var (
	ErrMissingMagic = errors.New("ar: invalid archive: missing magic number")
)

// Writer errors.
var (
	ErrFileNotDone  = errors.New("ar: file not fully written")
	ErrWriteTooLong = errors.New("ar: write too long")
)
