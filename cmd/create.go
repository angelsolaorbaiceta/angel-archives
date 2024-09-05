package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
	"github.com/dustin/go-humanize"
)

func CreateArchive(outFileName string, inFileNames []string) {
	fmt.Fprintf(os.Stderr, "Creating archive %s with %d files...\n", outFileName, len(inFileNames))

	archive, err := archive.Create(inFileNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating archive: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(outFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	err = archive.Write(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing archive: %v\n", err)
		os.Exit(1)
	}

	var (
		archSize   = humanize.Bytes(uint64(archive.TotalSize()))
		headerSize = humanize.Bytes(uint64(archive.Header.HeaderLength))
	)

	fmt.Fprintf(os.Stderr, "Archive created successfully.\n")
	fmt.Fprintf(os.Stderr, "	> Archive size = %s.\n", archSize)
	fmt.Fprintf(os.Stderr, "	> Header size = %s.\n", headerSize)
}
