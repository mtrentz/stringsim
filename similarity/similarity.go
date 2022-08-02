package similarity

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/antzucaro/matchr"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Similarity struct {
	Metric string  `json:"metric"`
	S1     string  `json:"s1"`
	S2     string  `json:"s2"`
	Score  float64 `json:"score"`
}

func printResults(similarities []Similarity) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Metric\tS1\tS2\tScore\n")
	for _, similarity := range similarities {
		fmt.Fprintf(w, "%s\t%s\t%s\t%f\n", similarity.Metric, similarity.S1, similarity.S2, similarity.Score)
	}
	w.Flush()
}

// Func that receives the metric name and return a function that
// receives two strings and returns a score as float64.
func GetSimilarityFunc(metric string) func(string, string) float64 {
	switch metric {
	case "jaro":
		return matchr.Jaro
	case "levenshtein":
		// Wrap the function to return the result as float64
		return func(s1 string, s2 string) float64 {
			return float64(matchr.Levenshtein(s1, s2))
		}
	case "dameraulevenshtein":
		// Wrap the function to return the result as float64
		return func(s1 string, s2 string) float64 {
			return float64(matchr.DamerauLevenshtein(s1, s2))
		}
	case "hamming":
		// Wrap the function to return the result as float64
		// and handle error
		return func(s1 string, s2 string) float64 {
			score, err := matchr.Hamming(s1, s2)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return float64(score)
		}
	default:
		fmt.Println("Metric not supported")
		os.Exit(1)
		return nil
	}
}

// Flow for calculating the similarities, printing results and
// exporting output when the amount of calculations is not too
// high. Output is sorted alphabetically, can be printed to stdout
// and is written all at once to a file.
func NormalFlow(mainStrings []string, subSlices [][]string, metric string, amountGoroutines int, StringFlags map[string]string, BoolFlags map[string]bool) {
	calculateSimilarity := GetSimilarityFunc(metric)

	var similarities []Similarity
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Add the amount of goroutines to the wait group
	wg.Add(amountGoroutines)

	// Start the goroutines
	// by looping over each subslice
	for _, subSlice := range subSlices {
		// Create a goroutine for each subslice
		go func(subSlice []string) {
			// Calculate the similarities for each subslice
			for _, s1 := range mainStrings {
				for _, s2 := range subSlice {
					// Calculate the similarity
					score := calculateSimilarity(s1, s2)
					// Create a new similarity object
					similarity := Similarity{
						Metric: cases.Title(language.Und, cases.NoLower).String(metric),
						S1:     s1,
						S2:     s2,
						Score:  score,
					}
					// Add the similarity to the slice
					mu.Lock()
					similarities = append(similarities, similarity)
					mu.Unlock()
				}
			}
			// Done with this goroutine
			wg.Done()
		}(subSlice)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Sort the slice by s1, alphabetically
	sort.Slice(similarities, func(i, j int) bool {
		return strings.ToLower(similarities[i].S1) < strings.ToLower(similarities[j].S1)
	})

	// Now check if its not set to silent to print results
	if !BoolFlags["Silent"] {
		printResults(similarities)
	}

	// Check if output to write to file
	if StringFlags["Output"] != "" {
		// Write to file
		writeToFile(StringFlags["Output"], similarities)
	}
}

// Flow for big files, for which I will not hold
// the similarities slice in memory and I'll
// be apending each result to the final json file.
func BigFileFlow(mainStrings []string, subSlices [][]string, metric string, amountGoroutines int, StringFlags map[string]string, BoolFlags map[string]bool) {
	calculateSimilarity := GetSimilarityFunc(metric)

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Add the amount of goroutines to the wait group
	wg.Add(amountGoroutines)

	// Create a json file with an empty array or a csv file,
	// depending on the extension. Will exit if the extension
	// is not supported.
	createEmptyFile(StringFlags["Output"])

	// Open the file
	file, err := os.OpenFile(StringFlags["Output"], os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// For a json output, the file will start as [] and I'll be appending
	// each result to the file.
	// It's important to note that the first time I'm appending
	// I'll have to ommit a comma, since the normal apending is
	// `,{"key":value}]`.
	isEmpty := isEmptyList(file)

	// Start the goroutines
	// by looping over each subslice
	for _, subSlice := range subSlices {
		// Create a goroutine for each subslice
		go func(subSlice []string) {
			// Calculate the similarities for each subslice
			for _, s1 := range mainStrings {
				for _, s2 := range subSlice {
					// Calculate the similarity
					score := calculateSimilarity(s1, s2)
					// Create a new similarity object
					similarity := Similarity{
						Metric: cases.Title(language.Und, cases.NoLower).String(metric),
						S1:     s1,
						S2:     s2,
						Score:  score,
					}
					// Lock the array and append the similarity
					mu.Lock()
					appendToFile(file, &similarity, isEmpty)
					if isEmpty {
						isEmpty = false
					}
					mu.Unlock()
				}
			}
			// Done with this goroutine
			wg.Done()
		}(subSlice)
	}

	// Wait for all goroutines to finish
	wg.Wait()

}
