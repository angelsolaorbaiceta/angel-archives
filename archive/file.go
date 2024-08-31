package archive

import (
	"io"
)

// ArchiveFile represents a single file in the archive.
type ArchiveFile struct {
	FileName        string
	CompressedBytes []byte
}

// CompressedSize returns the size of the compressed file in bytes.
func (f *ArchiveFile) CompressedSize() int32 {
	return int32(len(f.CompressedBytes))
}

// DecompressedBytes returns the uncompressed bytes of the file.
func (f *ArchiveFile) DecompressedBytes() ([]byte, error) {
	return Decompress(f.CompressedBytes)
}

// NewFileFromReader creates a new ArchiveFile from a reader.
// It reads its bytes, compresses them using gzip, and returns the ArchiveFile.
func NewFileFromReader(reader io.Reader, fileName string) (*ArchiveFile, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	compressedData, err := Compress(data)
	if err != nil {
		return nil, err
	}

	return &ArchiveFile{
		FileName:        fileName,
		CompressedBytes: compressedData,
	}, nil
}
