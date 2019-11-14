// Package ar implements access to ar archives.
//
// The ar file format that this package implements is the version from Plan 9.
// Modern versions like GNU or BSD are not supported.
package ar

const (
	// ARMAG is the ar magic number at the beginning of every file.
	ARMAG = "!<arch>\n"
	// ARFMAG is a magic number at the end of each file header.
	ARFMAG = "`\n"
)
