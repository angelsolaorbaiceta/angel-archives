package archive

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// magic is a unique identifier for the archive format.
// It's the ASCII representation of "AAR?".
var magic = []byte{0x41, 0x41, 0x52, 0x3F}

// byteOrder is the byte order used to serialize integers.
var byteOrder = binary.LittleEndian

// Header represents the metadata of the archive.
// It includes the header's length in bytes and a list of file entries.
type Header struct {
	// HeaderLength is the length of the header in bytes, including the magic and header length fields.
	HeaderLength uint32
	Entries      []*HeaderFileEntry
}

// Write writes the header into the provided writer.
//
// The header is serialized as follows:
//
//  1. The magic field is serialized as a 4-byte sequence.
//  2. The header length field is serialized as a 4-byte sequence.
//  3. Each file entry is serialized as follows:
//     - The first 2 bytes represent the length of the file name in bytes.
//     - The file name is serialized as a sequence of bytes.
//     - The offset field is serialized as a 4-byte sequence.
//     - The size field is serialized as a 4-byte sequence.
func (h *Header) Write(w io.Writer) error {
	bytesWritten := uint32(0)

	// Write the magic (4 bytes)
	w.Write(magic)
	bytesWritten += 4

	// Write the header length (4 bytes)
	if err := binary.Write(w, byteOrder, h.HeaderLength); err != nil {
		return err
	} else {
		bytesWritten += 4
	}

	for _, entry := range h.Entries {
		if err := entry.Write(w); err != nil {
			return err
		} else {
			bytesWritten += entry.totalBytes()
		}
	}

	// Check that the passed in header length matches the actual length of the
	// serialized header
	if bytesWritten != h.HeaderLength {
		return fmt.Errorf(
			"header length mismatch: expected %d, got %d", h.HeaderLength, bytesWritten,
		)
	}

	return nil
}

// ReadHeader reads the header from the provided reader and returns a Header struct.
func ReadHeader(r io.Reader) (*Header, error) {
	var (
		readMagic    = make([]byte, 4)
		headerLength uint32
		readBytes    uint32 = 0
		fileEntries  []*HeaderFileEntry
	)

	// Read 	// Read the magic (4 bytes)
	if _, err := io.ReadFull(r, readMagic); err != nil {
		return nil, err
	} else {
		readBytes += 4
	}

	// Check if the magic is correct
	if !bytes.Equal(magic, readMagic) {
		return nil, fmt.Errorf("invalid magic: got %v, expected %v", readMagic, magic)
	}

	// Read the header length (4 bytes)
	if err := binary.Read(r, byteOrder, &headerLength); err != nil {
		return nil, err
	} else {
		readBytes += 4
	}

	for readBytes < headerLength {
		entry, err := ReadHeaderFile(r)
		if err != nil {
			return nil, err
		} else {
			readBytes += entry.totalBytes()
		}

		fileEntries = append(fileEntries, entry)
	}

	return &Header{
		HeaderLength: headerLength,
		Entries:      fileEntries,
	}, nil
}
