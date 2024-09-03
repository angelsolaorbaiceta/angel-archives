package archive

import (
	"encoding/binary"
	"fmt"
	"io"
)

// magic is a unique identifier for the archive format.
// It's the ASCII representation of "AAR?".
var magic = []byte{0x41, 0x41, 0x52, 0x3F}

// magicLen is the length of the magic field in bytes.
const magicLen = uint32(4)

// byteOrder is the byte order used to serialize integers.
var byteOrder = binary.LittleEndian

// A Header represents the metadata of the archive.
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
// It doesn't close the reader.
func ReadHeader(r io.Reader) (*Header, error) {
	var (
		headerLength uint32
		readBytes    uint32 = 0
		fileEntries  []*HeaderFileEntry
	)

	if err := mustReadMagic(r); err != nil {
		return nil, err
	} else {
		readBytes += magicLen
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

// ErrEntryNotFoundInHeader is returned when a file entry is not found in the header.
var ErrEntryNotFoundInHeader = fmt.Errorf("entry not found in header")

// FindHeaderEntryByName uses the reader to read the header until a file with the
// provided name is found. It returns the file entry or a errEntryNotFoundInHeader
// error if the file is not found. Other errors can be returned if the reader fails.
// The reader isn't closed.
func FindHeaderEntryByName(r io.Reader, fileName string) (*HeaderFileEntry, error) {
	var (
		headerLength uint32
		readBytes    uint32 = 0
	)

	if err := mustReadMagic(r); err != nil {
		return nil, err
	} else {
		readBytes += magicLen
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

		if entry.Name == fileName {
			return entry, nil
		}
	}

	return nil, ErrEntryNotFoundInHeader
}
