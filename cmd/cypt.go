package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
)

func EncryptArchive(fileName, password string) {
	// Read the archive
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

	// Encrypt the archive
	encArch, err := arch.Encrypt(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting archive: %v\n", err)
		os.Exit(1)
	}

	// Write the encrypted archive to disk
	encFileName := fileName + ".enc"
	encFile, err := os.Create(encFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating encrypted archive file: %v\n", err)
		os.Exit(1)
	}
	defer encFile.Close()

	if err := encArch.Write(encFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing encrypted archive file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Archive encrypted successfully to %s\n", encFileName)

	// Remove the original archive
	if err := os.Remove(fileName); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing original archive: %v\n", err)
		os.Exit(1)
	}
}

func DecryptArchive(fileName, password string) {
	// Read the encrypted archive
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening encrypted archive file: %v\n", err)
		os.Exit(1)
	}

	encArch, err := archive.ReadEncryptedArchive(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading encrypted archive: %v\n", err)
		os.Exit(1)
	}

	// Decrypt the archive
	arch, err := encArch.Decrypt(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting archive: %v\n", err)
		os.Exit(1)
	}

	// Write the decrypted archive to disk
	decFileName := decryptFileName(fileName)
	decFile, err := os.Create(decFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating decrypted archive file: %v\n", err)
		os.Exit(1)
	}
	defer decFile.Close()

	if err := arch.Write(decFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing decrypted archive file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Archive decrypted successfully to %s\n", decFileName)

	// Remove the encrypted archive
	if err := os.Remove(fileName); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing encrypted archive: %v\n", err)
		os.Exit(1)
	}
}

// decryptFileName returns the decrypted file name from the encrypted file name.
// If the file name doesn't end with ".enc", it appends ".dec" to the file name.
// Otherwise, it removes the ".enc" extension.
func decryptFileName(fileName string) string {
	if len(fileName) < 4 || fileName[len(fileName)-4:] != ".enc" {
		return fileName + ".dec"
	}

	return fileName[:len(fileName)-4]
}
