package utils

import (
	"bufio"
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

// Reads strings from a txt file separated by newline
// or a json file as an array of strings.
func ReadFromFile(filename string) []string {
	// Read from txt file
	if ext := filepath.Ext(filename); ext == ".txt" {
		return ReadFromTxtFile(filename)
	}

	// Read from json file
	if ext := filepath.Ext(filename); ext == ".json" {
		return ReadFromJsonFile(filename)
	}

	// If file extension not txt or json
	// prints error and exits.
	fmt.Println("File extension not .txt or .json")
	os.Exit(1)
	return nil
}

// Reads all lines from a txt file and returns them as a list.
func ReadFromTxtFile(filename string) []string {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Close file
	file.Close()

	return lines
}

// Read from a json file that is a list of strings and returns them as a list.
func ReadFromJsonFile(filename string) []string {
	// Expecting a json file with a top level list of only strings
	var arr []string

	// Read content from files and unmarshal into struct
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&arr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return arr
}
