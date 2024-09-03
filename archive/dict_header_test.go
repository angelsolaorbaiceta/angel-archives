package archive

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDictHeader(t *testing.T) {
	data := []byte{
		0x41, 0x41, 0x52, 0x3F, // magic
		0x1A, 0x00, 0x00, 0x00, // header length
		0x08, 0x00, // length of file name
		't', 'e', 's', 't', '.', 't', 'x', 't', // file name
		0x1B, 0x00, 0x00, 0x00, // offset
		0x04, 0x00, 0x00, 0x00, // size
	}
	reader := bytes.NewReader(data)

	header, err := ReadDictHeader(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, header.HeaderLength, uint32(26))
	assert.Equal(t, len(header.Entries), 1)
	assert.Equal(t, header.Entries["test.txt"].Offset, uint32(27))
	assert.Equal(t, header.Entries["test.txt"].Size, uint32(4))
}
