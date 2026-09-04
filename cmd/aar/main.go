package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/angelsolaorbaiceta/aar/cmd"
)

func main() {
	var (
		createCmd                = flag.NewFlagSet("create", flag.ExitOnError)
		createFileNameFlag       = createCmd.String("f", "", "Output filename of the archive")
		createEncryptFlag        = createCmd.Bool("encrypt", false, "Encrypt the archive with a password")
		createDeleteOriginalFlag = createCmd.Bool("delete", false, "Delete the original files")
		createHelpFlag           = createCmd.Bool("help", false, "Show help for create command")

		extractCmd          = flag.NewFlagSet("extract", flag.ExitOnError)
		extractFileNameFlag = extractCmd.String("f", "", "Filename of the archive to extract")
		extractNameFlag     = extractCmd.String("n", "", "Extract a specific file by name from the archive")
		extractDecryptFlag  = extractCmd.Bool("decrypt", false, "Decrypt the archive before extracting")
		extractHelpFlag     = extractCmd.Bool("help", false, "Show help for extract command")

		listCmd          = flag.NewFlagSet("list", flag.ExitOnError)
		listFileNameFlag = listCmd.String("f", "", "Filename of the archive to list")
		listHelpFlag     = listCmd.Bool("help", false, "Show help for list command")

		encryptCmd          = flag.NewFlagSet("encrypt", flag.ExitOnError)
		encryptFileNameFlag = encryptCmd.String("f", "", "Filename of the archive to encrypt")
		encryptHelpFlag     = encryptCmd.Bool("help", false, "Show help for encrypt command")

		decryptCmd          = flag.NewFlagSet("decrypt", flag.ExitOnError)
		decryptFileNameFlag = decryptCmd.String("f", "", "Filename of the archive to decrypt")
		decryptHelpFlag     = decryptCmd.Bool("help", false, "Show help for decrypt command")
	)

	if len(os.Args) < 2 {
		showMainHelp()
		os.Exit(1)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		showMainHelp()
		return
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createHelpFlag {
			showCreateHelp()
			return
		}
		validateFileName(*createFileNameFlag)
		fileNames := createCmd.Args()
		createArchive(
			*createFileNameFlag,
			fileNames,
			*createEncryptFlag,
			*createDeleteOriginalFlag,
		)

	case "extract":
		extractCmd.Parse(os.Args[2:])
		if *extractHelpFlag {
			showExtractHelp()
			return
		}
		validateFileName(*extractFileNameFlag)

		if *extractNameFlag == "" {
			cmd.ExtractArchive(*extractFileNameFlag, *extractDecryptFlag)
		} else {
			cmd.ExtractArchiveFile(
				*extractFileNameFlag,
				*extractNameFlag,
				*extractDecryptFlag,
			)
		}

	case "list":
		listCmd.Parse(os.Args[2:])
		if *listHelpFlag {
			showListHelp()
			return
		}
		validateFileName(*listFileNameFlag)
		cmd.ListArchive(*listFileNameFlag)

	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		if *encryptHelpFlag {
			showEncryptHelp()
			return
		}
		validateFileName(*encryptFileNameFlag)
		password := cmd.PromptPasswordWithConfirmation()

		cmd.EncryptArchive(*encryptFileNameFlag, password)

	case "decrypt":
		decryptCmd.Parse(os.Args[2:])
		if *decryptHelpFlag {
			showDecryptHelp()
			return
		}
		validateFileName(*decryptFileNameFlag)
		password := cmd.PromptPassword()

		cmd.DecryptArchive(*decryptFileNameFlag, password)

	default:
		showMainHelp()
		os.Exit(1)
	}
}

func validateFileName(name string) {
	if name == "" {
		fmt.Fprintf(os.Stderr, "You must specify a filename with the -f flag.\n")
		os.Exit(1)
	}
}

func createArchive(fileName string, fileNames []string, encrypt, deleteOriginal bool) {
	if len(fileNames) == 0 {
		fmt.Fprintf(os.Stderr, "You must specify at least one file to add to the archive.\n")
		os.Exit(1)
	}

	// Add appropriate extension based on encryption flag
	finalFileName := addArchiveExtension(fileName, encrypt)

	cmd.CreateArchive(finalFileName, fileNames, encrypt, deleteOriginal)
}

func addArchiveExtension(fileName string, encrypt bool) string {
	// Remove existing .aar or .aar.enc extensions if present
	if strings.HasSuffix(fileName, ".aar.enc") {
		fileName = strings.TrimSuffix(fileName, ".aar.enc")
	} else if strings.HasSuffix(fileName, ".aar") {
		fileName = strings.TrimSuffix(fileName, ".aar")
	}

	// Add appropriate extension
	if encrypt {
		return fileName + ".aar.enc"
	}
	return fileName + ".aar"
}

func showMainHelp() {
	fmt.Fprintf(os.Stderr, "aar - Angel Archives\n\n")
	fmt.Fprintf(os.Stderr, "An archiving tool that xz-compresses and bundles files together into an archive.\n")
	fmt.Fprintf(os.Stderr, "Archives can be encrypted and decrypted for maximum privacy.\n\n")
	fmt.Fprintf(os.Stderr, "Usage: aar <command> [options]\n\n")
	fmt.Fprintf(os.Stderr, "Available commands:\n")
	fmt.Fprintf(os.Stderr, "  create   Create a new archive from files\n")
	fmt.Fprintf(os.Stderr, "  extract  Extract files from an archive\n")
	fmt.Fprintf(os.Stderr, "  list     List the contents of an archive\n")
	fmt.Fprintf(os.Stderr, "  encrypt  Encrypt an archive with a password\n")
	fmt.Fprintf(os.Stderr, "  decrypt  Decrypt an encrypted archive\n\n")
	fmt.Fprintf(os.Stderr, "Use 'aar <command> --help' for more information about a command.\n")
}

func showCreateHelp() {
	fmt.Fprintf(os.Stderr, "aar create - Create a new archive from files\n\n")
	fmt.Fprintf(os.Stderr, "Usage: aar create -f <archive> [--encrypt] <file1> [file2] ...\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -f <archive>  Output filename (extensions added automatically)\n")
	fmt.Fprintf(os.Stderr, "  --encrypt     Encrypt the archive with a password\n")
	fmt.Fprintf(os.Stderr, "  --help        Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "File extensions:\n")
	fmt.Fprintf(os.Stderr, "  Without --encrypt: .aar extension is added automatically\n")
	fmt.Fprintf(os.Stderr, "  With --encrypt:    .aar.enc extension is added automatically\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  aar create -f archive file1.txt file2.txt    # Creates archive.aar\n")
	fmt.Fprintf(os.Stderr, "  aar create -f backup *.txt                   # Creates backup.aar\n")
	fmt.Fprintf(os.Stderr, "  aar create -f project src/ docs/ README.md   # Creates project.aar\n")
	fmt.Fprintf(os.Stderr, "  aar create -f secret --encrypt file1.txt     # Creates secret.aar.enc\n")
}

func showExtractHelp() {
	fmt.Fprintf(os.Stderr, "aar extract - Extract files from an archive\n\n")
	fmt.Fprintf(os.Stderr, "Usage: aar extract -f <archive> [--decrypt] [-n <filename>]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -f <archive>   Filename of the archive to extract\n")
	fmt.Fprintf(os.Stderr, "  -n <filename>  Extract a specific file by name (optional)\n")
	fmt.Fprintf(os.Stderr, "  --decrypt      Decrypt the archive before extracting\n")
	fmt.Fprintf(os.Stderr, "  --help         Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  aar extract -f archive.aar\n")
	fmt.Fprintf(os.Stderr, "  aar extract -f archive.aar -n file2.txt\n")
	fmt.Fprintf(os.Stderr, "  aar extract -f secret.aar.enc --decrypt\n")
	fmt.Fprintf(os.Stderr, "  aar extract -f secret.aar.enc --decrypt -n file2.txt\n")
}

func showListHelp() {
	fmt.Fprintf(os.Stderr, "aar list - List the contents of an archive\n\n")
	fmt.Fprintf(os.Stderr, "Usage: aar list -f <archive>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -f <archive>  Filename of the archive to list\n")
	fmt.Fprintf(os.Stderr, "  --help        Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  aar list -f archive.aarch\n")
	fmt.Fprintf(os.Stderr, "  aar list -f backup.aarch\n")
}

func showEncryptHelp() {
	fmt.Fprintf(os.Stderr, "aar encrypt - Encrypt an archive with a password\n\n")
	fmt.Fprintf(os.Stderr, "Usage: aar encrypt -f <archive>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -f <archive>  Filename of the archive to encrypt\n")
	fmt.Fprintf(os.Stderr, "  --help        Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Description:\n")
	fmt.Fprintf(os.Stderr, "  Encrypts an archive using AES-256-GCM algorithm. The password must be\n")
	fmt.Fprintf(os.Stderr, "  at least 8 characters long. Creates a new file with .enc extension and\n")
	fmt.Fprintf(os.Stderr, "  removes the original archive.\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  aar encrypt -f archive.aarch\n")
	fmt.Fprintf(os.Stderr, "  aar encrypt -f backup.aarch\n")
}

func showDecryptHelp() {
	fmt.Fprintf(os.Stderr, "aar decrypt - Decrypt an encrypted archive\n\n")
	fmt.Fprintf(os.Stderr, "Usage: aar decrypt -f <encrypted_archive>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -f <encrypted_archive>  Filename of the encrypted archive to decrypt\n")
	fmt.Fprintf(os.Stderr, "  --help                  Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Description:\n")
	fmt.Fprintf(os.Stderr, "  Decrypts an archive using AES-256-GCM algorithm. Creates a new file\n")
	fmt.Fprintf(os.Stderr, "  without the .enc extension and removes the encrypted archive.\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  aar decrypt -f archive.aarch.enc\n")
	fmt.Fprintf(os.Stderr, "  aar decrypt -f backup.aarch.enc\n")
}
