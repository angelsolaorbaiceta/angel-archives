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

// FileEntry represents a single file's metadata in the archive.
type FileEntry struct {
	// Name is a unique identifier for the file.
	Name string
	// Offset is the byte offset from the beginning of the archive where the file's data begins.
	// Uses 4 bytes to store the offset.
	Offset uint32
	// Size is the size of the file's data in bytes. Uses 4 bytes to store the size.
	Size uint32
}

// nameLength returns the length of the file name in bytes.
func (f *FileEntry) nameLength() uint16 {
	return uint16(len(f.Name))
}

// Header represents the metadata of the archive.
// It includes the header's length in bytes and a list of file entries.
type Header struct {
	// HeaderLength is the length of the header in bytes, including the magic and header length fields.
	HeaderLength uint32
	Entries      []FileEntry
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

	// Write the magic
	buffer.Write(magic)

	// Write the header length
	if err := binary.Write(buffer, byteOrder, h.HeaderLength); err != nil {
		return nil, err
	}

	for _, entry := range h.Entries {
		// Write the length of the file name in bytes
		if err := binary.Write(buffer, byteOrder, entry.nameLength()); err != nil {
			return nil, err
		}

		// Write the file name
		if _, err := buffer.WriteString(entry.Name); err != nil {
			return nil, err
		}

		// Write the offset
		if err := binary.Write(buffer, byteOrder, entry.Offset); err != nil {
			return nil, err
		}

		// Write the size
		if err := binary.Write(buffer, byteOrder, entry.Size); err != nil {
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
	header := &Header{}

	// Read the magic
	readMagic := make([]byte, 4)
	if _, err := io.ReadFull(r, readMagic); err != nil {
		return nil, err
	}

	// Check if the magic is correct
	if !bytes.Equal(magic, readMagic) {
		return nil, fmt.Errorf("invalid magic: got %v, expected %v", readMagic, magic)
	}

	// Read the header length
	if err := binary.Read(r, byteOrder, &header.HeaderLength); err != nil {
		return nil, err
	}

	for {
		entry := FileEntry{}

		// Read the length of the file name
		var length uint16
		if err := binary.Read(r, byteOrder, &length); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Read the file name
		name := make([]byte, length)
		if _, err := io.ReadFull(r, name); err != nil {
			return nil, err
		}
		entry.Name = string(name)

		// Read the offset
		if err := binary.Read(r, byteOrder, &entry.Offset); err != nil {
			return nil, err
		}

		// Read the size
		if err := binary.Read(r, byteOrder, &entry.Size); err != nil {
			return nil, err
		}

		header.Entries = append(header.Entries, entry)
	}

	return header, nil
}
