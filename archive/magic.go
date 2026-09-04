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
)

// magic is a unique identifier for the archive format.
// It's the ASCII representation of "AAR?".
var magic = []byte{0x41, 0x41, 0x52, 0x3F}

// encMagic is a unique identifier for the encrypted archive format.
// It's the ASCII representation of "AARX".
var encMagic = []byte{0x41, 0x41, 0x52, 0x58}

// magicLen is the length of the magic field in bytes.
const magicLen = uint32(4)

// ErrInvalidMagic is returned when the magic field is not correct.
var ErrInvalidMagic = fmt.Errorf("invalid magic, expected %v", magic)

// ErrInvalidEncMagic is returned when the magic field is not correct.
var ErrInvalidEncMagic = fmt.Errorf("invalid magic, expected %v", encMagic)

// readMagic identifies the magic bytes in the reader, and returns its type.
func ReadMagic(r io.Reader) FileMagicType {
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

// mustReadMagic reads the magic field from the provided reader.
// If the magic field is not correct, it returns an error.
func mustReadMagic(r io.Reader) error {
	readMagic := make([]byte, 4)

	if _, err := io.ReadFull(r, readMagic); err != nil {
		return err
	}

	if !bytes.Equal(magic, readMagic) {
		return ErrInvalidMagic
	}

	return nil
}

// mustReadEncryptedMagic reads the magic field from the provided reader.
func mustReadEncryptedMagic(r io.Reader) error {
	readMagic := make([]byte, 4)

	if _, err := io.ReadFull(r, readMagic); err != nil {
		return err
	}

	if !bytes.Equal(encMagic, readMagic) {
		return ErrInvalidEncMagic
	}

	return nil
}
