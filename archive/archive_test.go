package archive

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArchive(t *testing.T) {
	var (
		fileOne      = createTempFileForTest(t, "fileOne.txt", "AAAAAAAA")
		fileTwo      = createTempFileForTest(t, "fileTwo.txt", "BBBBBBBB")
		archive, err = CreateArchive([]string{fileOne.FileName, fileTwo.FileName})

		wantHeaderLen = uint32(8 + (2 + len(fileOne.FileName) + 8) + (2 + len(fileTwo.FileName) + 8))
	)

	assert.Nil(t, err)

	t.Run("archive header length", func(t *testing.T) {
		assert.Equal(t, wantHeaderLen, archive.Header.HeaderLength)
	})

	t.Run("archive header first file entry", func(t *testing.T) {
		got := archive.Header.Entries[0]
		want := &HeaderFileEntry{
			Name:   fileOne.FileName,
			Offset: wantHeaderLen + 1,
			Size:   fileOne.CompressedSize(),
		}

		assert.Equal(t, want, got)
	})

	t.Run("archive header second file entry", func(t *testing.T) {
		got := archive.Header.Entries[1]
		want := &HeaderFileEntry{
			Name:   fileTwo.FileName,
			Offset: wantHeaderLen + 1 + fileOne.CompressedSize(),
			Size:   fileTwo.CompressedSize(),
		}

		assert.Equal(t, want, got)
	})

	t.Run("archive files", func(t *testing.T) {
		assert.Equal(t, 2, len(archive.Files))

		t.Run("first file", func(t *testing.T) {
			got, err := archive.Files[0].DecompressedBytes()
			assert.Nil(t, err)
			assert.Equal(t, []byte("AAAAAAAA"), got)
		})

		t.Run("second file", func(t *testing.T) {
			got, err := archive.Files[1].DecompressedBytes()
			assert.Nil(t, err)
			assert.Equal(t, []byte("BBBBBBBB"), got)
		})
	})
}

func TestWriteAndReadArchive(t *testing.T) {
	var (
		fileOne    = createTempFileForTest(t, "fileOne.txt", "AAAAAAAA")
		fileTwo    = createTempFileForTest(t, "fileTwo.txt", "BBBBBBBB")
		archive, _ = CreateArchive([]string{fileOne.FileName, fileTwo.FileName})
		writer     = new(bytes.Buffer)
	)

	assert.Nil(t, archive.Write(writer))

	reader := bytes.NewReader(writer.Bytes())
	got, err := ReadArchive(reader)

	assert.Nil(t, err)
	assert.Equal(t, archive, got)
}

// createTempFileForTest creates a file in the test's temporal directory and returns
// an ArchiveFile with the file's name and expected compressed bytes.
func createTempFileForTest(t *testing.T, fileName, content string) *ArchiveFile {
	filePath := filepath.Join(t.TempDir(), fileName)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}

	compressedBytes, err := Compress([]byte(content))
	if err != nil {
		t.Fatalf("Error compressing file: %v", err)
	}

	return &ArchiveFile{
		FileName:        filePath,
		CompressedBytes: compressedBytes,
	}
}
