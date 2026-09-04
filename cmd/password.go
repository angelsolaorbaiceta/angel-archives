package cmd

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func promptPasswordWithConfirmation() string {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening terminal: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	fmt.Fprint(tty, "Password: ")
	passwordBytes, err := term.ReadPassword(int(tty.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Move to the next line after password input
	fmt.Fprintln(tty)

	password := string(passwordBytes)
	validatePassword(password)

	fmt.Fprint(tty, "Confirm password: ")
	passwordConfirmationBytes, err := term.ReadPassword(int(tty.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Move to the next line after password input
	fmt.Fprintln(tty)

	passwordConfirmation := string(passwordConfirmationBytes)
	if password != passwordConfirmation {
		fmt.Fprintf(os.Stderr, "Passwords do not match.\n")
		os.Exit(1)
	}

	return password
}

func promptPassword() string {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening terminal: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	fmt.Fprint(tty, "Password: ")
	passwordBytes, err := term.ReadPassword(int(tty.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Move to the next line after password input
	fmt.Fprintln(tty)

	return string(passwordBytes)
}

func validatePassword(password string) {
	if len(password) < 8 {
		fmt.Fprintf(os.Stderr, "The password must be at least 8 characters long.\n")
		os.Exit(1)
	}
}
