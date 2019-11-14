package ar

import (
	"os"
	"time"
)

// Header is a file header within an ar archive.
type Header struct {
	Name string
	Date time.Time
	UID  uint32
	GID  uint32
	Mode os.FileMode
	Size int
}
