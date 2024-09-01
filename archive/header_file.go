package archive

import (
	"encoding/binary"
	"io"
)

// HeaderFileEntry represents a single file's metadata in the archive.
type HeaderFileEntry struct {
	// Name is a unique identifier for the file.
	Name string
	// Offset is the byte offset from the beginning of the archive where the file's data begins.
	// Uses 4 bytes to store the offset.
	Offset uint32
	// Size is the size of the file's data in bytes. Uses 4 bytes to store the size.
	Size uint32
}

// nameLength returns the length of the file name in bytes.
func (f *HeaderFileEntry) nameLength() uint16 {
	return uint16(len(f.Name))
}

// totalBytes returns the total number of bytes required to serialize the HeaderFileEntry.
// This includes the length of the file name (2 bytes), the file name itself, the
// offset (4 bytes), and the size (4 bytes).
func (f *HeaderFileEntry) totalBytes() uint32 {
	return 2 + uint32(f.nameLength()) + 4 + 4
}

// Write writes the serialized HeaderFileEntry to the provided writer.
func (f *HeaderFileEntry) Write(w io.Writer) error {
	// Write the length of the file name in bytes (2 bytes)
	if err := binary.Write(w, byteOrder, f.nameLength()); err != nil {
		return err
	}

	// Write the file name (nameLength bytes)
	if _, err := w.Write([]byte(f.Name)); err != nil {
		return err
	}

	// Write the offset (4 bytes)
	if err := binary.Write(w, byteOrder, f.Offset); err != nil {
		return err
	}

	// Write the size (4 bytes)
	if err := binary.Write(w, byteOrder, f.Size); err != nil {
		return err
	}

	return nil
}

// ReadHeaderFile reads a HeaderFileEntry from the provided reader.
func ReadHeaderFile(r io.Reader) (*HeaderFileEntry, error) {
	var (
		nameLength uint16
		name       []byte
		offset     uint32
		size       uint32
	)

	// Read the file name length (2 bytes)
	if err := binary.Read(r, byteOrder, &nameLength); err != nil {
		return nil, err
	}

	// Read the file name (nameLength bytes)
	name = make([]byte, nameLength)
	if _, err := io.ReadFull(r, name); err != nil {
		return nil, err
	}

	// Read the offset (4 bytes)
	if err := binary.Read(r, byteOrder, &offset); err != nil {
		return nil, err
	}

	// Read the size (4 bytes)
	if err := binary.Read(r, byteOrder, &size); err != nil {
		return nil, err
	}

	return &HeaderFileEntry{
		Name:   string(name),
		Offset: offset,
		Size:   size,
	}, nil
}
