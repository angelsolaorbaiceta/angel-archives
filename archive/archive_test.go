package archive

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArchive(t *testing.T) {
	var (
		fileOne      = createTempFileForTest(t, "fileOne.txt", "AAAAAAAA")
		fileTwo      = createTempFileForTest(t, "fileTwo.txt", "BBBBBBBB")
		archive, err = Create([]string{fileOne.FileName, fileTwo.FileName}, nil)

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
		archive, _ = Create([]string{fileOne.FileName, fileTwo.FileName}, nil)
		writer     = new(bytes.Buffer)
	)

	assert.Nil(t, archive.Write(writer))

	reader := bytes.NewReader(writer.Bytes())
	got, err := ReadArchive(reader)

	assert.Nil(t, err)
	assert.Equal(t, archive, got)
}

func TestReadFileByName(t *testing.T) {
	var (
		fileOne    = createTempFileForTest(t, "fileOne.txt", "AAAAAAAA")
		fileTwo    = createTempFileForTest(t, "fileTwo.txt", "BBBBBBBB")
		archive, _ = Create([]string{fileOne.FileName, fileTwo.FileName}, nil)
		arBytes    = new(bytes.Buffer)
	)

	archive.Write(arBytes)

	t.Run("read fileOne", func(t *testing.T) {
		var (
			reader = bytes.NewReader(arBytes.Bytes())
			got, _ = ReadFileByName(reader, fileOne.FileName)
		)

		assert.Equal(t, fileOne, got)
	})

	t.Run("read fileTwo", func(t *testing.T) {
		var (
			reader = bytes.NewReader(arBytes.Bytes())
			got, _ = ReadFileByName(reader, fileTwo.FileName)
		)

		assert.Equal(t, fileTwo, got)
	})
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

	return NewFileFromCompressedBytes(filePath, compressedBytes)
}

func TestCreateArchiveReportsProgress(t *testing.T) {
	var (
		fileOne = createTempFileForTest(t, "fileOne.txt", "AAAAAAAA")
		fileTwo = createTempFileForTest(t, "fileTwo.txt", "BBBBBBBB")
		paths   = []string{fileOne.FileName, fileTwo.FileName}
		got     []Progress
	)

	_, err := Create(paths, func(p Progress) { got = append(got, p) })
	assert.Nil(t, err)

	t.Run("reports every file once", func(t *testing.T) {
		assert.Equal(t, len(paths), len(got))

		reported := make([]string, 0, len(got))
		for _, p := range got {
			reported = append(reported, p.Path)
		}
		assert.ElementsMatch(t, paths, reported)
	})

	t.Run("counts up to the total", func(t *testing.T) {
		for i, p := range got {
			assert.Equal(t, i+1, p.Done)
			assert.Equal(t, len(paths), p.Total)
			assert.Nil(t, p.Err)
			assert.NotZero(t, p.Compressed)
		}
	})
}

func TestCreateArchiveReportsFailures(t *testing.T) {
	var (
		fileOne = createTempFileForTest(t, "fileOne.txt", "AAAAAAAA")
		missing = filepath.Join(t.TempDir(), "missing.txt")
		got     []Progress
	)

	_, err := Create(
		[]string{fileOne.FileName, missing},
		func(p Progress) { got = append(got, p) },
	)

	t.Run("returns the error", func(t *testing.T) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("still reports both files", func(t *testing.T) {
		assert.Equal(t, 2, len(got))

		var failed int
		for _, p := range got {
			if p.Err != nil {
				failed++
				assert.Equal(t, missing, p.Path)
			}
		}
		assert.Equal(t, 1, failed)
	})
}
