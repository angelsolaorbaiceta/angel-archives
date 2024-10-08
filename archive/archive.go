package archive

import (
	"bytes"
	"io"
)

// An Archive represents a collection of files stored in a single file.
type Archive struct {
	Header *Header
	Files  []*ArchiveFile
}

// TotalSize returns the total size of the archive in bytes.
// It includes the header and all the files' compressed data.
func (a *Archive) TotalSize() uint64 {
	var total uint64 = uint64(a.Header.HeaderLength)

	for _, file := range a.Files {
		total += uint64(file.CompressedSize())
	}

	return total
}

// GetBytes returns the archive as a byte slice.
func (a *Archive) GetBytes() ([]byte, error) {
	data := new(bytes.Buffer)
	if err := a.Write(data); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

// Write writes the archive into the provided writer.
func (a *Archive) Write(w io.Writer) error {
	if err := a.Header.Write(w); err != nil {
		return err
	}

	for _, file := range a.Files {
		if err := file.Write(w); err != nil {
			return err
		}
	}

	return nil
}

// ReadArchive reads an archive from the provided reader.
// It reads all the files and the header, and returns an Archive struct.
// It doesn't close the reader.
func ReadArchive(r io.Reader) (*Archive, error) {
	header, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}

	files, err := ReadFiles(r, header)
	if err != nil {
		return nil, err
	}

	return &Archive{
		Header: header,
		Files:  files,
	}, nil
}

// Create creates a new archive from the provided file paths.
func Create(filePaths []string) (*Archive, error) {
	files, err := readFiles(filePaths)
	if err != nil {
		return nil, err
	}

	header, err := makeHeader(files)
	if err != nil {
		return nil, err
	}

	return &Archive{
		Header: header,
		Files:  files,
	}, nil
}

// readFiles reads the files concurrently from the provided file paths.
// Each file is xz-compressed and stored in an ArchiveFile struct.
// The order of the files is preserved.
func readFiles(filePaths []string) ([]*ArchiveFile, error) {
	type item struct {
		file *ArchiveFile
		err  error
		idx  int
	}

	var (
		files = make([]*ArchiveFile, len(filePaths))
		ch    = make(chan item, len(filePaths))
	)

	for i, path := range filePaths {
		go func(path string) {
			file, err := NewFileFromPath(path)
			ch <- item{file, err, i}
		}(path)
	}

	for range filePaths {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}

		files[it.idx] = it.file
	}

	return files, nil
}

func makeHeader(files []*ArchiveFile) (*Header, error) {
	var (
		entries           = make([]*HeaderFileEntry, len(files))
		totalBytes uint32 = 8
	)

	for i, file := range files {
		entries[i] = NewHeaderFileEntry(file.FileName, file.CompressedSize())
		totalBytes += entries[i].totalBytes()
	}

	currentOffset := totalBytes + 1
	for _, entry := range entries {
		entry.Offset = currentOffset
		currentOffset += entry.Size
	}

	return &Header{
		HeaderLength: totalBytes,
		Entries:      entries,
	}, nil
}

// ReadFileByName reads the archive's header until the name of the file is found.
// Then, it reads the file's data and returns an ArchiveFile struct.
// If the file is not found, it returns an ErrEntryNotFoundInHeader error.
func ReadFileByName(r ReaderSeeker, fileName string) (*ArchiveFile, error) {
	if fileHeaderEntry, err := FindHeaderEntryByName(r, fileName); err != nil {
		return nil, err
	} else {
		return fileHeaderEntry.ReadFrom(r)
	}
}
