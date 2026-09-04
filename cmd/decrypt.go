package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
	"github.com/spf13/cobra"
)

var (
	decryptFileName string

	decryptCmd = &cobra.Command{
		Use:                   "decrypt -f <encrypted_archive>",
		Short:                 "Decrypt an encrypted archive",
		DisableFlagsInUseLine: true,
		Long: `Decrypt an encrypted archive.

Decrypts an archive using AES-256-GCM algorithm. Creates a new file
without the .enc extension and removes the encrypted archive.`,
		Example: `  aar decrypt -f archive.aar.enc
  aar decrypt -f backup.aar.enc`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			password := promptPassword()
			decryptArchive(decryptFileName, password)
		},
	}
)

func init() {
	addFileNameFlag(decryptCmd, &decryptFileName, "Filename of the encrypted archive to decrypt")

	rootCmd.AddCommand(decryptCmd)
}

func decryptArchive(fileName, password string) {
	// Read the encrypted archive
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening encrypted archive file: %v\n", err)
		os.Exit(1)
	}

	// The whole archive is read into memory, so the file can be closed right
	// away. It needs to be closed before removing it below.
	encArch, err := archive.ReadEncryptedArchive(reader)
	reader.Close()
	if err != nil {
		if err == archive.ErrWrongFileType {
			fmt.Fprintf(os.Stderr, "Can only decrypt .aar.enc files\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error reading encrypted archive: %v\n", err)
		}
		os.Exit(1)
	}

	// Decrypt the archive
	arch, err := encArch.Decrypt(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting archive: %v\n", err)
		os.Exit(1)
	}

	// Write the decrypted archive to disk
	decFileName := decryptedFileName(fileName)
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

// decryptedFileName returns the decrypted file name from the encrypted file name.
// If the file name doesn't end with ".enc", it appends ".dec" to the file name.
// Otherwise, it removes the ".enc" extension.
func decryptedFileName(fileName string) string {
	if len(fileName) < 4 || fileName[len(fileName)-4:] != ".enc" {
		return fileName + ".dec"
	}

	return fileName[:len(fileName)-4]
}
