package archive

import (
	"bytes"
	"fmt"
	"io"
)

type FileMagicType string

const (
	FileMagicUnknown    FileMagicType = "Unknown"
	FileMagicArchive    FileMagicType = "Archive"
	FileMagicEncArchive FileMagicType = "EncArchive"

	// magicLen is the length of the magic identifier.
	magicLen uint32 = 4
)

var (
	// magic is a unique identifier for the archive format.
	// It's the ASCII representation of "AAR?".
	magic = []byte{0x41, 0x41, 0x52, 0x3F}

	// encMagic is a unique identifier for the encrypted archive format.
	// It's the ASCII representation of "AARX".
	encMagic = []byte{0x41, 0x41, 0x52, 0x58}

	// ErrWrongFileType is the error returned when the magic doesn't match with
	// what's expected in a specific context.
	ErrWrongFileType = fmt.Errorf("wrong file type")
)

// readMagic identifies the magic bytes in the reader, and returns its type.
func readMagic(r io.Reader) FileMagicType {
	readMagic := make([]byte, 4)
	if _, err := io.ReadFull(r, readMagic); err != nil {
		return FileMagicUnknown
	}

	if bytes.Equal(magic, readMagic) {
		return FileMagicArchive
	}

	if bytes.Equal(encMagic, readMagic) {
		return FileMagicEncArchive
	}

	return FileMagicUnknown
}
