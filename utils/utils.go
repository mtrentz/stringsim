package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mtrentz/similarity-cli/similarity"
	"github.com/spf13/cobra"
)

// Check for a minimum amount of arguments, if not enough,
// prints help page, error and exits.
func CheckForMinimumArgs(cmd *cobra.Command, n int, args []string) {
	if len(args) < n {
		cmd.Usage()
		fmt.Printf("Expected %d arguments, got %d.\n", n, len(args))
		os.Exit(1)
	}
}

// Writes a list of similarities to a json file as a list.
func WriteToFile(filename string, similarities []similarity.Similarity) {
	// Check if extension is already .json, else add it
	if ext := filepath.Ext(filename); ext != ".json" {
		filename = filename + ".json"
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	j, err := json.MarshalIndent(similarities, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write to file
	file.Write(j)
}
