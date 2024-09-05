package archive

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptArchive(t *testing.T) {
	archive := makeTestArchive()

	encrypted, err := archive.Encrypt("password")
	assert.Nil(t, err)

	decrypted, err := encrypted.Decrypt("password")
	assert.Nil(t, err)

	assert.Equal(t, archive, decrypted)
}

func TestWriteAndReadEncryptedArchive(t *testing.T) {
	var (
		archive      = makeTestArchive()
		encrypted, _ = archive.Encrypt("password")
		w            = new(bytes.Buffer)
	)

	err := encrypted.Write(w)
	assert.Nil(t, err)

	r := bytes.NewReader(w.Bytes())
	readArchive, err := ReadEncryptedArchive(r)
	assert.Nil(t, err)

	assert.Equal(t, encrypted, readArchive)
}

func makeTestArchive() *Archive {
	return &Archive{
		Header: &Header{
			HeaderLength: 46,
			Entries: []*HeaderFileEntry{
				{
					Name:   "file1.txt",
					Size:   12,
					Offset: 28,
				},
				{
					Name:   "file2.txt",
					Size:   16,
					Offset: 40,
				},
			},
		},
		Files: []*ArchiveFile{
			{
				FileName: "file1.txt",
				CompressedBytes: []byte{
					0x78, 0x9c, 0x4b, 0x4c,
					0x4f, 0x49, 0x2d, 0x2e,
					0x01, 0x00, 0x00, 0xff,
				},
			},
			{
				FileName: "file2.txt",
				CompressedBytes: []byte{
					0x78, 0x9c, 0x4b, 0x4c,
					0x4f, 0x49, 0x2d, 0x2e,
					0x01, 0x00, 0x00, 0xff,
					0x78, 0x9c, 0x4b, 0x4c,
				},
			},
		},
	}
}
