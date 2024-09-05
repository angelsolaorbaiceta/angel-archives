package cmd

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func PromptPassword() string {
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Move to the next line after password input
	fmt.Println()

	password := string(passwordBytes)
	validatePassword(password)

	fmt.Print("Confirm password: ")
	passwordConfirmationBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Move to the next line after password input
	fmt.Println()

	passwordConfirmation := string(passwordConfirmationBytes)
	if password != passwordConfirmation {
		fmt.Fprintf(os.Stderr, "Passwords do not match.\n")
		os.Exit(1)
	}

	return password
}

func validatePassword(password string) {
	if len(password) < 8 {
		fmt.Fprintf(os.Stderr, "The password must be at least 8 characters long.\n")
		os.Exit(1)
	}
}
