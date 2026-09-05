package archive

import (
	"io"
	"os"
)

// ArchiveFile represents a single file in the archive.
// It includes the file's name and its raw bytes.
type ArchiveFile struct {
	FileName string
	Data     []byte
}

// Write writes the bytes of the file into the provided writer.
func (f *ArchiveFile) Write(w io.Writer) error {
	_, err := w.Write(f.Data)
	return err
}

// SizeInBytes returns the number of bytes in the file.
func (f *ArchiveFile) SizeInBytes() int {
	return len(f.Data)
}

// NewFileFromReader creates a new ArchiveFile from the bytes read from the reader.
func NewFileFromReader(reader io.Reader, fileName string) (*ArchiveFile, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return NewFileFromData(data, fileName), nil
}

// NewFileFromData returns a new file from the raw bytes and file name.
func NewFileFromData(data []byte, fileName string) *ArchiveFile {
	return &ArchiveFile{
		FileName: fileName,
		Data:     data,
	}
}

// NewFileFromPath creates a new ArchiveFile from tye bytes in the file at the file path.
// Returns an error if the file path can't be opened or the file can't be read from.
func NewFileFromPath(path string) (*ArchiveFile, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return NewFileFromReader(reader, path)
}

// ReadFiles reads the files sequentially from the provided reader using the header.
func ReadFiles(r io.Reader, header *Header) ([]*ArchiveFile, error) {
	files := make([]*ArchiveFile, len(header.Entries))

	for i, entry := range header.Entries {
		fileData := make([]byte, entry.Size)
		if _, err := r.Read(fileData); err != nil {
			return nil, err
		}

		files[i] = &ArchiveFile{
			FileName: entry.Name,
			Data:     fileData,
		}
	}

	return files, nil
}
