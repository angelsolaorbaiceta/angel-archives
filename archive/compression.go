package archive

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz"
)

// Compress compresses the given bytes using the xz algorithm and returns the
// compressed bytes.
func Compress(data []byte) ([]byte, error) {
	var (
		compressedData bytes.Buffer
		xzWriter, err  = xz.NewWriter(&compressedData)
	)

	if err != nil {
		return nil, err
	}

	// Write the data to the xz writer
	if _, err := xzWriter.Write(data); err != nil {
		return nil, err
	}

	// Close the xz writer to flush the compressed data
	if err := xzWriter.Close(); err != nil {
		return nil, err
	}

	return compressedData.Bytes(), nil
}

// Decompress decompresses the given bytes using the xz algorithm and returns the
// uncompressed bytes.
func Decompress(data []byte) ([]byte, error) {
	var (
		reader        = bytes.NewReader(data)
		xzReader, err = xz.NewReader(reader)
	)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(xzReader)
}
