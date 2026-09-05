package archive

import (
	"bytes"
	"errors"
	"io"
)

// An Archive represents a collection of xz-compressed files stored in a single file.
type Archive struct {
	Header *Header
	Files  []*ArchiveFile
}

// TotalSize returns the total size of the archive in bytes.
// It includes the header and all the files' raw data.
func (a *Archive) TotalSize() uint64 {
	var total uint64 = uint64(a.Header.HeaderLength)

	for _, file := range a.Files {
		total += uint64(file.SizeInBytes())
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
//
// It expects to read from an unencrypted archive, returns an error otherwise.
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
// If onProgress is non-nil, it's called once per file as it finishes being read
// and compressed.
func Create(filePaths []string) (*Archive, error) {
	files, err := mapFilesToArchives(filePaths)
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

// mapFilesToArchives reads each of the file paths content and creates an ArchiveFile
// for each of them.
func mapFilesToArchives(filePaths []string) ([]*ArchiveFile, error) {
	var (
		files = make([]*ArchiveFile, len(filePaths))
		errs  = make([]error, len(filePaths))
	)

	for i, path := range filePaths {
		files[i], errs[i] = NewFileFromPath(path)
	}

	return files, errors.Join(errs...)
}

func makeHeader(files []*ArchiveFile) (*Header, error) {
	var (
		entries           = make([]*HeaderFileEntry, len(files))
		totalBytes uint32 = 8
	)

	for i, file := range files {
		entries[i] = NewHeaderFileEntry(file.FileName, uint32(file.SizeInBytes()))
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
