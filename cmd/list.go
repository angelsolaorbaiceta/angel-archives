package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
)

func ListArchive(fileName string) {
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive file: %v\n", err)
		os.Exit(1)
	}

	header, err := archive.ReadHeader(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading archive header: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Archive has the following files:\n")
	for _, entry := range header.Entries {
		fmt.Fprintf(os.Stdout, "	> %s\n", entry)
	}
}
