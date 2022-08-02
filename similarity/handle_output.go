package similarity

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Detect if output is to json or csv, write it all at once,
// which works for the smaller files that everything is hold in memory.
func writeToFile(filename string, similarities []Similarity) {
	// Check the extension
	ext := filepath.Ext(filename)

	// If the extension is .json, write to json
	if ext == ".json" {
		writeToJson(filename, similarities)
		return
	}

	// If the extension is .csv, write to csv
	if ext == ".csv" {
		writeToCsv(filename, similarities)
		return
	}

	// If the extension is neither .json nor .csv,
	// prints error and exit
	fmt.Println("File extension is not .json or .csv")
	os.Exit(1)
}

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
	defer file.Close()

	j, err := json.MarshalIndent(similarities, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write to file
	file.Write(j)
}

// Write file to CSV all at once including headers.
func writeToCsv(filename string, similarities []Similarity) {
	// Check if extension is already .csv, else add it
	if ext := filepath.Ext(filename); ext != ".csv" {
		filename = filename + ".csv"
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	// Write header
	csvWriter.Write([]string{"metric", "s1", "s2", "score"})

	// Write similarities
	for _, similarity := range similarities {
		csvWriter.Write([]string{similarity.Metric, similarity.S1, similarity.S2, fmt.Sprintf("%f", similarity.Score)})
	}

	csvWriter.Flush()
}

// By the filename, either create and empty csv
// with the headers, or an empty json array.
func createEmptyFile(filename string) {
	// Check the extension
	ext := filepath.Ext(filename)

	// If the extension is .json, create an empty json array
	if ext == ".json" {
		createEmptyJsonArrayFile(filename)
		return
	}

	// If the extension is .csv, create an empty csv with the headers
	if ext == ".csv" {
		createEmptyCsvFile(filename)
		return
	}

	// If the extension is neither .json nor .csv,
	// prints error and exit
	fmt.Println("File extension is not .json or .csv")
	os.Exit(1)
}

// Create and empty csv with the headers
func createEmptyCsvFile(filename string) {
	// Check if extension is already .csv, else add it
	if ext := filepath.Ext(filename); ext != ".csv" {
		filename = filename + ".csv"
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Write header
	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"metric", "s1", "s2", "score"})
	csvWriter.Flush()
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
	defer file.Close()

	// Write to file
	file.Write([]byte("[]\n"))

	// Close file
	file.Close()
}

// Append to the correct file, depending on the file extension.
func appendToFile(file *os.File, similarity *Similarity, isEmpty bool) {

	extension := filepath.Ext(file.Name())

	// If the extension is .json, append to json
	if extension == ".json" {
		appendToJsonArray(file, similarity, isEmpty)
		return
	}

	// If the extension is .csv, append to csv
	if extension == ".csv" {
		appendToCsv(file, similarity)
		return
	}

	// If the extension is neither .json nor .csv,
	// prints error and exit
	fmt.Println("Error appending to file. File extension is not .json or .csv")
	os.Exit(1)
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

// Append line at the end of a csv file.
func appendToCsv(file *os.File, similarity *Similarity) {
	file.Seek(0, 2)
	w := csv.NewWriter(file)
	w.Write([]string{similarity.Metric, similarity.S1, similarity.S2, fmt.Sprintf("%f", similarity.Score)})
	w.Flush()
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
