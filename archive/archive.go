package archive

type Archive struct {
	Header *Header
	Files  []*ArchiveFile
}

// CreateArchive creates a new archive from the provided file paths.
func CreateArchive(filePaths []string) (*Archive, error) {
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

func readFiles(filePaths []string) ([]*ArchiveFile, error) {
	files := make([]*ArchiveFile, len(filePaths))

	for i, path := range filePaths {
		file, err := NewFileFromPath(path)
		if err != nil {
			return nil, err
		}

		files[i] = file
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
