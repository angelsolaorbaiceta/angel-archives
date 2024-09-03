package archive

import "io"

// ReaderSeeker is an interface that combines the io.Reader and io.Seeker interfaces.
type ReaderSeeker interface {
	io.Reader
	io.Seeker
}
