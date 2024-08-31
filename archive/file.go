package archive

import (
	"bytes"
	"compress/gzip"
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

// UncompressedBytes returns the uncompressed bytes of the file.
func (f *ArchiveFile) UncompressedBytes() ([]byte, error) {
	var (
		reader          = bytes.NewReader(f.CompressedBytes)
		gzipReader, err = gzip.NewReader(reader)
	)

	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return io.ReadAll(gzipReader)
}

// NewFileFromReader creates a new ArchiveFile from a reader.
// It reads its bytes, compresses them using gzip, and returns the ArchiveFile.
func NewFileFromReader(reader io.Reader, fileName string) (*ArchiveFile, error) {
	// Read the bytes from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var (
		compressedData bytes.Buffer
		gzipWriter     = gzip.NewWriter(&compressedData)
	)

	// Write the data to the gzip writer
	if _, err := gzipWriter.Write(data); err != nil {
		return nil, err
	}

	// Close the gzip writer to flush the compressed data
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return &ArchiveFile{
		FileName:        fileName,
		CompressedBytes: compressedData.Bytes(),
	}, nil
}
