package archive

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileFromReader(t *testing.T) {
	var (
		fileName = "test.txt"
		data     = []byte("AAAAAAAA")
		reader   = bytes.NewReader(data)
	)

	archiveFile, err := NewFileFromReader(reader, fileName)

	assert.Nil(t, err)
	assert.Equal(t, archiveFile.FileName, fileName)
	assert.NotEmpty(t, archiveFile.CompressedBytes)

	uncompressed, _ := archiveFile.UncompressedBytes()
	assert.Equal(t, data, uncompressed)
}
