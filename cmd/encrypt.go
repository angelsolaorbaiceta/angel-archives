package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
	"github.com/spf13/cobra"
)

var (
	encryptFileName string

	encryptCmd = &cobra.Command{
		Use:                   "encrypt -f <archive>",
		Short:                 "Encrypt an archive with a password",
		DisableFlagsInUseLine: true,
		Long: `Encrypt an archive with a password.

Encrypts an archive using AES-256-GCM algorithm. The password must be
at least 8 characters long. Creates a new file with .enc extension and
removes the original archive.`,
		Example: `  aar encrypt -f archive.aar
  aar encrypt -f backup.aar`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			password := promptPasswordWithConfirmation()
			encryptArchive(encryptFileName, password)
		},
	}
)

func init() {
	addFileNameFlag(encryptCmd, &encryptFileName, "Filename of the archive to encrypt")

	rootCmd.AddCommand(encryptCmd)
}

func encryptArchive(fileName, password string) {
	// Read the archive
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive file: %v\n", err)
		os.Exit(1)
	}

	// The whole archive is read into memory, so the file can be closed right
	// away. It needs to be closed before removing it below.
	arch, err := archive.ReadArchive(reader)
	reader.Close()
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
