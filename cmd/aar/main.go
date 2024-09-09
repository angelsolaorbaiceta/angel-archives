package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/cmd"
)

func main() {
	var (
		createCmd          = flag.NewFlagSet("create", flag.ExitOnError)
		createFileNameFlag = createCmd.String("f", "", "Output filename of the archive")

		extractCmd          = flag.NewFlagSet("extract", flag.ExitOnError)
		extractFileNameFlag = extractCmd.String("f", "", "Filename of the archive to extract")
		extractNameFlag     = extractCmd.String("n", "", "Extract a specific file by name from the archive")

		listCmd          = flag.NewFlagSet("list", flag.ExitOnError)
		listFileNameFlag = listCmd.String("f", "", "Filename of the archive to list")

		encryptCmd          = flag.NewFlagSet("encrypt", flag.ExitOnError)
		encryptFileNameFlag = encryptCmd.String("f", "", "Filename of the archive to encrypt")

		decryptCmd          = flag.NewFlagSet("decrypt", flag.ExitOnError)
		decryptFileNameFlag = decryptCmd.String("f", "", "Filename of the archive to decrypt")
	)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: aar <command> [options]\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		validateFileName(*createFileNameFlag)
		fileNames := createCmd.Args()
		createArchive(*createFileNameFlag, fileNames)

	case "extract":
		extractCmd.Parse(os.Args[2:])
		validateFileName(*extractFileNameFlag)

		if *extractNameFlag == "" {
			cmd.ExtractArchive(*extractFileNameFlag)
		} else {
			cmd.ExtractArchiveFile(*extractFileNameFlag, *extractNameFlag)
		}

	case "list":
		listCmd.Parse(os.Args[2:])
		validateFileName(*listFileNameFlag)
		cmd.ListArchive(*listFileNameFlag)

	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		validateFileName(*encryptFileNameFlag)
		password := cmd.PromptPasswordWithConfirmation()

		cmd.EncryptArchive(*encryptFileNameFlag, password)

	case "decrypt":
		decryptCmd.Parse(os.Args[2:])
		validateFileName(*decryptFileNameFlag)
		password := cmd.PromptPassword()

		cmd.DecryptArchive(*decryptFileNameFlag, password)

	default:
		fmt.Fprintf(os.Stderr, "Usage: aar <command> [options]\n")
		os.Exit(1)
	}
}

func validateFileName(name string) {
	if name == "" {
		fmt.Fprintf(os.Stderr, "You must specify a filename with the -f flag.\n")
		os.Exit(1)
	}
}

func createArchive(fileName string, fileNames []string) {
	if len(fileNames) == 0 {
		fmt.Fprintf(os.Stderr, "You must specify at least one file to add to the archive.\n")
		os.Exit(1)
	}

	cmd.CreateArchive(fileName, fileNames)
}
