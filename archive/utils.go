package archive

import (
	"bytes"
	"fmt"
	"io"
)

// mustReadMagic reads the magic field from the provided reader.
// If the magic field is not correct, it returns an error.
func mustReadMagic(r io.Reader) error {
	readMagic := make([]byte, 4)

	// Read the magic (4 bytes)
	if _, err := io.ReadFull(r, readMagic); err != nil {
		return err
	}

	// Check if the magic is correct
	if !bytes.Equal(magic, readMagic) {
		return fmt.Errorf("invalid magic: got %v, expected %v", readMagic, magic)
	}

	return nil
}
