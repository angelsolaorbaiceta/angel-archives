package archive

import (
	"bytes"
	"errors"
	"io"
	"runtime"
	"sync"
)

// An Archive represents a collection of xz-compressed files stored in a single file.
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
func Create(filePaths []string, onProgress func(Progress)) (*Archive, error) {
	files, err := readAndCompressFiles(filePaths, onProgress)
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

// Progress reports a single file that has finished being read and compressed.
// Err is non-nil if that file failed.
type Progress struct {
	Path       string
	Err        error
	Done       int // Files finished so far, including this one
	Total      int
	Compressed uint32 // Compressed size, 0 if the file failed
}

// readAndCompressFiles reads the files concurrently from the provided file paths.
// Each file is xz-compressed and stored in an ArchiveFile struct.
// The order of the files is preserved, and the work is bounded to GOMAXPROCS
// files in flight at a time.
//
// If onProgress is non-nil, it's called once per file as it finishes, in
// completion order. It runs on the caller's goroutine, so it doesn't need to be
// safe for concurrent use.
func readAndCompressFiles(
	filePaths []string,
	onProgress func(Progress),
) ([]*ArchiveFile, error) {
	var (
		files  = make([]*ArchiveFile, len(filePaths))
		errs   = make([]error, len(filePaths))
		sem    = make(chan struct{}, runtime.GOMAXPROCS(0))
		events = make(chan Progress)
		wg     sync.WaitGroup
	)

	// The files are dispatched from their own goroutine so that the caller is
	// free to consume events while files are still being compressed. Dispatching
	// them here would deadlock: the workers would block sending to events, never
	// releasing their semaphore slot, and this loop would stall on a full sem.
	go func() {
		for i, path := range filePaths {
			wg.Add(1)
			sem <- struct{}{} // Blocks once N are in flight

			go func(i int, path string) {
				defer wg.Done()
				defer func() { <-sem }()

				files[i], errs[i] = NewFileFromPath(path)

				event := Progress{Path: path, Err: errs[i], Total: len(filePaths)}
				if errs[i] == nil {
					event.Compressed = files[i].CompressedSize()
				}

				events <- event
			}(i, path)
		}

		wg.Wait()
		close(events)
	}()

	var done int
	for event := range events {
		done++
		event.Done = done

		if onProgress != nil {
			onProgress(event)
		}
	}

	return files, errors.Join(errs...)
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
