package archive

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteHeader(t *testing.T) {
	header := Header{
		// 4 bytes for magic +
		// 4 bytes for header length +
		// 2 bytes for file name length +
		// 8 bytes for file name +
		// 8 bytes for offset and size = 26 bytes
		HeaderLength: 26,
		Entries: []*HeaderFileEntry{
			{
				Name:   "test.txt",
				Offset: 27, // 26 bytes for the header + 1 byte
				Size:   4,
			},
		},
	}

	writer := new(bytes.Buffer)
	err := header.Write(writer)
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
	got := writer.Bytes()

	assert.Equal(t, len(got), len(expected))
	assert.Equal(t, got, expected)
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

func TestFindHeaderEntryByName(t *testing.T) {
	header := Header{
		HeaderLength: 45,
		Entries: []*HeaderFileEntry{
			{
				Name:   "test.txt",
				Offset: 27,
				Size:   4,
			},
			{
				Name:   "test2.txt",
				Offset: 31,
				Size:   5,
			},
		},
	}
	headerBytes := new(bytes.Buffer)
	header.Write(headerBytes)

	reader := bytes.NewReader(headerBytes.Bytes())
	entry, err := FindHeaderEntryByName(reader, "test.txt")

	assert.Nil(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, entry.Name, "test.txt")
	assert.Equal(t, entry.Offset, uint32(27))
	assert.Equal(t, entry.Size, uint32(4))
}

func TestReadFrom(t *testing.T) {
	var (
		data = []byte{
			0x41, 0x41, 0x52, 0x3F, // magic
			0x1A, 0x00, 0x00, 0x00, // header length
			0x08, 0x00, // length of file name
			't', 'e', 's', 't', '.', 't', 'x', 't', // file name
			0x1B, 0x00, 0x00, 0x00, // offset
			0x04, 0x00, 0x00, 0x00, // size
			0x41, 0x41, 0x41, 0x41, // file data
		}
		reader = bytes.NewReader(data)
		entry  = &HeaderFileEntry{
			Name:   "test.txt",
			Offset: 27,
			Size:   4,
		}
		want = NewFileFromCompressedBytes("test.txt", []byte("AAAA"))
	)

	file, err := entry.ReadFrom(reader)

	assert.Nil(t, err)
	assert.Equal(t, file, want)
}
