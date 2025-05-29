package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
)

func ExtractArchive(fileName string, decrypt bool) {
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive file: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	var arch *archive.Archive
	if decrypt {
		encryptedArch, err := archive.ReadEncryptedArchive(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading encrypted archive: %v\n", err)
			os.Exit(1)
		}
		
		password := PromptPassword()
		arch, err = encryptedArch.Decrypt(password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decrypting archive: %v\n", err)
			os.Exit(1)
		}
	} else {
		arch, err = archive.ReadArchive(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading archive: %v\n", err)
			os.Exit(1)
		}
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

func ExtractArchiveFile(fileName, fileToExtract string, decrypt bool) {
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive file: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	var archFile *archive.ArchiveFile
	if decrypt {
		encryptedArch, err := archive.ReadEncryptedArchive(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading encrypted archive: %v\n", err)
			os.Exit(1)
		}
		
		password := PromptPassword()
		arch, err := encryptedArch.Decrypt(password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decrypting archive: %v\n", err)
			os.Exit(1)
		}
		
		var found bool
		for _, file := range arch.Files {
			if file.FileName == fileToExtract {
				archFile = file
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "File not found in archive: %s\n", fileToExtract)
			fmt.Fprintf(os.Stderr, "Use the list command to see the files in the archive.\n")
			os.Exit(1)
		}
	} else {
		archFile, err = archive.ReadFileByName(reader, fileToExtract)
		if err != nil {
			if err == archive.ErrEntryNotFoundInHeader {
				fmt.Fprintf(os.Stderr, "File not found in archive: %s\n", fileToExtract)
				fmt.Fprintf(os.Stderr, "Use the list command to see the files in the archive.\n")
			} else {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			}

			os.Exit(1)
		}
	}

	outFile, err := os.Create(archFile.FileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	err = archFile.WriteDecompressed(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
}
