package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
)

func CreateArchive(outFileName string, inFileNames []string) {
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
}
