package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/aar/archive"
	"github.com/spf13/cobra"
)

var (
	listFileName string

	listCmd = &cobra.Command{
		Use:                   "list -f <archive>",
		Short:                 "List the contents of an archive",
		DisableFlagsInUseLine: true,
		Example: `  aar list -f archive.aar
  aar list -f backup.aar`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			listArchive(listFileName)
		},
	}
)

func init() {
	addFileNameFlag(listCmd, &listFileName, "Filename of the archive to list")

	rootCmd.AddCommand(listCmd)
}

func listArchive(fileName string) {
	reader, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive file: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

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
