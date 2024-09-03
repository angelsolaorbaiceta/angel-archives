package archive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressAndDecompress(t *testing.T) {
	data := []byte("hello world")

	compressed, err := Compress(data)
	assert.Nil(t, err)

	decompressed, err := Decompress(compressed)
	assert.Nil(t, err)

	assert.Equal(t, data, decompressed)
}
