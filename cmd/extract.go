package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
)

func ExtractArchive(fileName string) {
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive file: %v\n", err)
		os.Exit(1)
	}

	arch, err := archive.ReadArchive(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading archive: %v\n", err)
		os.Exit(1)
	}

	for _, file := range arch.Files {
		fmt.Fprintf(os.Stderr, "Extracting %s...\n", file.FileName)
		outFile, err := os.Create(file.FileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			os.Exit(1)
		}
		defer outFile.Close()

		err = file.WriteDecompressed(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
	}
}
