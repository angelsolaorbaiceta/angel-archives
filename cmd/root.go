package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aar",
	Short: "aar - Angel Archives",
	Long: `aar - Angel Archives

An archiving tool that xz-compresses and bundles files together into an archive.
Archives can be encrypted and decrypted for maximum privacy.`,
	// With no subcommand there's nothing to do: show the help and exit with an
	// error code, as the previous flag-based implementation did.
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// Execute runs the root command, which dispatches to the appropriate
// subcommand based on the command line arguments.
func Execute() {
	// Cobra already reports the error and the usage, so there's nothing left
	// to print here.
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// addFileNameFlag adds the required "-f" flag to a command, binding it to the
// given target, and documenting it with the given usage string.
func addFileNameFlag(cmd *cobra.Command, target *string, usage string) {
	cmd.Flags().StringVarP(target, "file", "f", "", usage)
	cmd.MarkFlagRequired("file")
}
