package archive

import (
	"encoding/binary"
	"fmt"
	"io"
)

// HeaderFileEntry represents a single file's metadata in the archive.
type HeaderFileEntry struct {
	// Name is a unique identifier for the file.
	Name string
	// Offset is the byte offset from the beginning of the archive where the file's data begins.
	// Uses 4 bytes to store the offset.
	Offset uint32
	// Size is the size of the file's compressed data in bytes. Uses 4 bytes to store the size.
	Size uint32
}

// NewHeaderFileEntry creates a new header file entry with the given name and size,
// setting the offset at 0, as it's impossible to know the offset until the whole
// file is set up.
func NewHeaderFileEntry(name string, size uint32) *HeaderFileEntry {
	return &HeaderFileEntry{
		Name:   name,
		Offset: 0,
		Size:   size,
	}
}

// String returns a string representation of the HeaderFileEntry.
func (f *HeaderFileEntry) String() string {
	return fmt.Sprintf("%s (Offset in file: %d bytes, Compressed size: %d bytes)", f.Name, f.Offset, f.Size)
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

// ReadFrom reads the file data from the provided ReaderSeeker, using the file's
// offset and size
func (f *HeaderFileEntry) ReadFrom(r ReaderSeeker) (*ArchiveFile, error) {
	fileData := make([]byte, f.Size)

	if _, err := r.Seek(int64(f.Offset-1), io.SeekStart); err != nil {
		return nil, err
	}

	if _, err := r.Read(fileData); err != nil {
		return nil, err
	}

	return NewFileFromCompressedBytes(f.Name, fileData), nil
}
