package archive

import (
	"encoding/binary"
	"io"
)

// A DictHeader is the header of the file where the file entries are stored in a
// dictionary format, by name. Use this version to extract files by name.
// This header shouldn't be used for writing to the archive file.
type DictHeader struct {
	HeaderLength uint32
	Entries      map[string]*HeaderFileEntry
}

func ReadDictHeader(r io.Reader) (*DictHeader, error) {
	var (
		headerLength uint32
		readBytes    uint32 = 0
		fileEntries         = make(map[string]*HeaderFileEntry)
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
			fileEntries[entry.Name] = entry
		}
	}

	return &DictHeader{
		HeaderLength: headerLength,
		Entries:      fileEntries,
	}, nil
}
