package main

import "os"

func main() {
	if err := New().Execute(); err != nil {
		// Cobra will already print an error on problems like an unknown command
		// or missing required flag. Set an exit status of 1 on error, but don't
		// print it again.
		os.Exit(1)
	}
}
