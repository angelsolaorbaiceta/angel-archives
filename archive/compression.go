package archive

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Compress compresses the given bytes using gzip and returns the compressed
// bytes.
func Compress(data []byte) ([]byte, error) {
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

	return compressedData.Bytes(), nil
}

// Decompress decompresses the given bytes using gzip and returns the
// uncompressed bytes.
func Decompress(data []byte) ([]byte, error) {
	var (
		reader          = bytes.NewReader(data)
		gzipReader, err = gzip.NewReader(reader)
	)

	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return io.ReadAll(gzipReader)
}
