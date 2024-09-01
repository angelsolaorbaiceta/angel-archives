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
var byteOrder = binary.LittleEndian

// Header represents the metadata of the archive.
// It includes the header's length in bytes and a list of file entries.
type Header struct {
	// HeaderLength is the length of the header in bytes, including the magic and header length fields.
	HeaderLength uint32
	Entries      []*HeaderFileEntry
}

// Serialize serializes the header into a byte slice.
//
// The header is serialized as follows:
//
// 1. The magic field is serialized as a 4-byte sequence.
// 2. The header length field is serialized as a 4-byte sequence.
// 3. Each file entry is serialized as follows:
//   - The first 2 bytes represent the length of the file name in bytes.
//   - The file name is serialized as a sequence of bytes.
//   - The offset field is serialized as a 4-byte sequence.
//   - The size field is serialized as a 4-byte sequence.
func (h *Header) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// Write the magic (4 bytes)
	buffer.Write(magic)

	// Write the header length (4 bytes)
	if err := binary.Write(buffer, byteOrder, h.HeaderLength); err != nil {
		return nil, err
	}

	for _, entry := range h.Entries {
		if err := entry.Serialize(buffer); err != nil {
			return nil, err
		}
	}

	// Check that the passed in header length matches the actual length of the serialized header
	if uint32(buffer.Len()) != h.HeaderLength {
		return nil, fmt.Errorf("header length mismatch: expected %d, got %d", h.HeaderLength, buffer.Len())
	}

	return buffer.Bytes(), nil
}

// WriteHeader writes the serialized header to the provided writer.
func (h *Header) WriteHeader(w io.Writer) error {
	data, err := h.Serialize()
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
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
		entry, err := DeserializeHeaderFile(r)
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
