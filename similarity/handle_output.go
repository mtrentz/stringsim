package similarity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Writes a list of similarities to a json file as a list.
func writeToJson(filename string, similarities []Similarity) {
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

// Create an empty output json file with an empty array
func createEmptyJsonArrayFile(filename string) {
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

	// Write to file
	file.Write([]byte("[]\n"))

	// Close file
	file.Close()
}

// Appens a new object to the end of a json list.
// This only works for a list, since it seeks the end of file,
// works backwards until the last ']', and then add a new object
// at the end.
func appendToJsonArray(file *os.File, similarity *Similarity, isEmpty bool) {
	// Go to the end of the file
	file.Seek(-1, 2)
	b := make([]byte, 1)

	var err error
	var i int64

	// Read the last 3 bytes looking for a ']'
	for i = 1; i <= 3; i++ {
		_, err = file.Read(b)
		file.Seek(-1-i, 2)
		if err != nil {
			fmt.Println(err)
		}
		s := string(b)
		if s == "]" {
			file.Seek(-i, 2)
			j, _ := json.Marshal(similarity)
			if !isEmpty {
				file.Write([]byte(","))
			}
			file.Write(j)
			file.Write([]byte("]\n"))
			return
		}
	}

	// If we get here, the file is not a json list
	fmt.Println("File is not a json list")
	os.Exit(1)
}

// Detects if a file is "[]" or "[ ]"
func isEmptyList(file *os.File) bool {
	// Go to the beginning of the file
	file.Seek(0, 0)
	// Read first string
	firstChar := readByteAsString(file)

	// Advance 1 byte
	file.Seek(1, 0)
	secondChar := readByteAsString(file)

	// Go one more byte forward
	file.Seek(2, 0)
	thirdChar := readByteAsString(file)

	// Now check if the file is "[]" or "[ ]"
	if firstChar == "[" && secondChar == "]" {
		return true
	}
	if firstChar == "[" && secondChar == " " && thirdChar == "]" {
		return true
	}

	return false
}

func readByteAsString(file *os.File) string {
	b := make([]byte, 1)
	_, err := file.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
