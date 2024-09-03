package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/cmd"
)

var (
	createFlag      = flag.Bool("c", false, "Create a new archive")
	extractFlag     = flag.Bool("x", false, "Extract an archive")
	extractNameFlag = flag.String("n", "", "Extract a specific file by name from the archive")
	listFlag        = flag.Bool("l", false, "List the contents of an archive")
	fileName        = flag.String("f", "", "Output filename of the archive")
)

func main() {
	flag.Parse()
	validateFlags()

	if *createFlag {
		createArchive()
	} else if *extractFlag {
		if *extractNameFlag != "" {
			extractArchiveFile(*extractNameFlag)
		} else {
			extractArchive()
		}
	} else if *listFlag {
		listArchive()
	} else {
		fmt.Println("Usage: aar [options] -f <filename>")
		flag.PrintDefaults()
	}
}

func validateFlags() {
	activeFlags := 0
	if *createFlag {
		activeFlags++
	}
	if *extractFlag {
		activeFlags++
	}
	if *listFlag {
		activeFlags++
	}

	if activeFlags != 1 {
		fmt.Fprintf(os.Stderr, "You must specify one of the -c (create), -x (extract), or -l (list) flags.\n")
		os.Exit(1)
	}

	if *fileName == "" {
		fmt.Fprintf(os.Stderr, "You must specify a filename with the -f flag.\n")
		os.Exit(1)
	}
}

func createArchive() {
	fileNames := flag.Args()
	if len(fileNames) == 0 {
		fmt.Fprintf(os.Stderr, "You must specify at least one file to add to the archive.\n")
		os.Exit(1)
	}

	cmd.CreateArchive(*fileName, fileNames)
}

func extractArchive() {
	cmd.ExtractArchive(*fileName)
}

func extractArchiveFile(fileToExtract string) {
	cmd.ExtractArchiveFile(*fileName, fileToExtract)
}

func listArchive() {
	cmd.ListArchive(*fileName)
}
