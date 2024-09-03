package archive

import (
	"io"
	"os"
)

// ArchiveFile represents a single file in the archive.
// It includes the file's name and its compressed bytes (using xz).
// The decompressed bytes can be obtained using the DecompressedBytes method.
type ArchiveFile struct {
	FileName        string
	CompressedBytes []byte
}

// Write writes the compressed bytes of the file into the provided writer.
func (f *ArchiveFile) Write(w io.Writer) error {
	_, err := w.Write(f.CompressedBytes)
	return err
}

// WriteDecompressed writes the decompressed bytes of the file into the provided writer.
func (f *ArchiveFile) WriteDecompressed(w io.Writer) error {
	decompressedBytes, err := f.DecompressedBytes()
	if err != nil {
		return err
	}

	_, err = w.Write(decompressedBytes)
	return err
}

// CompressedSize returns the size of the compressed file in bytes.
func (f *ArchiveFile) CompressedSize() uint32 {
	return uint32(len(f.CompressedBytes))
}

// DecompressedBytes returns the uncompressed bytes of the file.
func (f *ArchiveFile) DecompressedBytes() ([]byte, error) {
	return Decompress(f.CompressedBytes)
}

// NewFileFromReader creates a new ArchiveFile from a reader.
// It reads its bytes, compresses them using xz, and returns the ArchiveFile.
func NewFileFromReader(reader io.Reader, fileName string) (*ArchiveFile, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	compressedData, err := Compress(data)
	if err != nil {
		return nil, err
	}

	return &ArchiveFile{
		FileName:        fileName,
		CompressedBytes: compressedData,
	}, nil
}

// NewFileFromPath creates a new ArchiveFile from a file path.
func NewFileFromPath(path string) (*ArchiveFile, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

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
			FileName:        entry.Name,
			CompressedBytes: fileData,
		}
	}

	return files, nil
}
