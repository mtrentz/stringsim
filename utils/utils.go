package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

// Reads strings from a txt file separated by newline
// or a json file as an array of strings.
func ReadFromFile(filename string) []string {
	// Read from txt file
	if ext := filepath.Ext(filename); ext == ".txt" {
		return readFromTxtFile(filename)
	}

	// Read from json file
	if ext := filepath.Ext(filename); ext == ".json" {
		return readFromJsonFile(filename)
	}

	// If file extension not txt or json
	// prints error and exits.
	fmt.Println("File extension not .txt or .json")
	os.Exit(1)
	return nil
}

// Reads all lines from a txt file and returns them as a list.
func readFromTxtFile(filename string) []string {
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
func readFromJsonFile(filename string) []string {
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

// Split slice into 'n' subslices as evenly as possible.
func SliceSplit(slice []string, n int) [][]string {
	// First create and slice with 'n' amount of slices
	subSlices := make([][]string, n)
	// Now loop over 'slice' and append each element into
	// one slice until there is no more elements left.
	sliceToAppend := 0
	for {
		if len(slice) == 0 {
			break
		}
		subSlices[sliceToAppend] = append(subSlices[sliceToAppend], slice[0])
		slice = slice[1:]
		sliceToAppend++
		if sliceToAppend == n {
			sliceToAppend = 0
		}
	}
	return subSlices
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
