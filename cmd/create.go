package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/angelsolaorbaiceta/aar/archive"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var (
	createFileName       string
	createEncrypt        bool
	createDeleteOriginal bool

	createCmd = &cobra.Command{
		Use:                   "create -f <archive> [--encrypt] <file1> [file2] ...",
		Short:                 "Create a new archive from files",
		DisableFlagsInUseLine: true,
		Long: `Create a new archive from files.

File extensions:
  Without --encrypt: .aar extension is added automatically
  With --encrypt:    .aar.enc extension is added automatically`,
		Example: `  aar create -f archive file1.txt file2.txt    # Creates archive.aar
  aar create -f backup *.txt                   # Creates backup.aar
  aar create -f project src/ docs/ README.md   # Creates project.aar
  aar create -f secret --encrypt file1.txt     # Creates secret.aar.enc`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Add appropriate extension based on encryption flag
			outFileName := addArchiveExtension(createFileName, createEncrypt)
			createArchive(outFileName, args, createEncrypt, createDeleteOriginal)
		},
	}
)

func init() {
	addFileNameFlag(createCmd, &createFileName, "Output filename of the archive")
	createCmd.Flags().BoolVar(&createEncrypt, "encrypt", false, "Encrypt the archive with a password")
	createCmd.Flags().BoolVar(&createDeleteOriginal, "delete", false, "Delete the original files")

	rootCmd.AddCommand(createCmd)
}

// addArchiveExtension returns the file name with the extension that corresponds
// to the archive type: ".aar" for plain archives and ".aar.enc" for encrypted
// ones. Any of those extensions already present in the name is replaced.
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

func createArchive(
	outFileName string,
	inFileNames []string,
	encrypt, deleteOriginal bool,
) {
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

	if encrypt {
		password := promptPasswordWithConfirmation()
		encryptedArchive, err := archive.Encrypt(password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting archive: %v\n", err)
			os.Exit(1)
		}

		err = encryptedArchive.Write(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing encrypted archive: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = archive.Write(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing archive: %v\n", err)
			os.Exit(1)
		}
	}

	if deleteOriginal {
		for _, filePath := range inFileNames {
			if err := os.Remove(filePath); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing file '%s': %v\n", filePath, err)
				continue
			}
		}
	}

	var (
		archSize   = humanize.Bytes(uint64(archive.TotalSize()))
		headerSize = humanize.Bytes(uint64(archive.Header.HeaderLength))
	)

	fmt.Fprintf(os.Stderr, "Archive created successfully.\n")
	fmt.Fprintf(os.Stderr, "	> Archive size = %s.\n", archSize)
	fmt.Fprintf(os.Stderr, "	> Header size = %s.\n", headerSize)
	fmt.Fprintf(os.Stderr, "Files in archive:\n")
	for _, file := range archive.Files {
		size := humanize.Bytes(uint64(file.CompressedSize()))
		fmt.Fprintf(os.Stderr, "	> %s (compressed size = %s)\n", file.FileName, size)
	}
}
