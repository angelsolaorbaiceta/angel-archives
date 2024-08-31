package archive

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderSerialization(t *testing.T) {
	header := Header{
		// 4 bytes for magic +
		// 4 bytes for header length +
		// 2 bytes for file name length +
		// 8 bytes for file name +
		// 8 bytes for offset and size = 26 bytes
		HeaderLength: 26,
		Entries: []HeaderFileEntry{
			{
				Name:   "test.txt",
				Offset: 27, // 26 bytes for the header + 1 byte
				Size:   4,
			},
		},
	}

	data, err := header.Serialize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte{
		0x41, 0x41, 0x52, 0x3F, // magic
		0x1A, 0x00, 0x00, 0x00, // header length
		0x08, 0x00, // length of file name
		't', 'e', 's', 't', '.', 't', 'x', 't', // file name
		0x1B, 0x00, 0x00, 0x00, // offset
		0x04, 0x00, 0x00, 0x00, // size
	}

	assert.Equal(t, len(data), len(expected))
	assert.Equal(t, data, expected)
}

func TestReadHeader(t *testing.T) {
	data := []byte{
		0x41, 0x41, 0x52, 0x3F, // magic
		0x1A, 0x00, 0x00, 0x00, // header length
		0x08, 0x00, // length of file name
		't', 'e', 's', 't', '.', 't', 'x', 't', // file name
		0x1B, 0x00, 0x00, 0x00, // offset
		0x04, 0x00, 0x00, 0x00, // size
	}
	reader := bytes.NewReader(data)

	header, err := ReadHeader(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, header.HeaderLength, uint32(26))
	assert.Equal(t, len(header.Entries), 1)
	assert.Equal(t, header.Entries[0].Name, "test.txt")
	assert.Equal(t, header.Entries[0].Offset, uint32(27))
	assert.Equal(t, header.Entries[0].Size, uint32(4))
}
