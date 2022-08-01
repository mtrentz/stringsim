package similarity

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

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

func (s *Similarity) Result() {
	fmt.Printf("Similarity between %s and %s using %s is %f\n", s.S1, s.S2, s.Metric, s.Score)
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
// high.
func NormalFlow(mainStrings []string, otherStrings []string, StringFlags map[string]string, BoolFlags map[string]bool) {
	// Getting the function based on the metric to use
	var Metric string
	var calculateSimilarity func(string, string) float64

	// Decide on the metric to use if has a flag
	if StringFlags["Metric"] != "" {
		// Make sure Metric is all lower case
		Metric = strings.ToLower(StringFlags["Metric"])
	} else {
		Metric = "jaro"
	}
	calculateSimilarity = GetSimilarityFunc(Metric)

	// Create a channel to guarantee max amount of goroutines
	// equal to cpu cores
	MAX_CPU_CORES := runtime.NumCPU()
	waitChan := make(chan struct{}, MAX_CPU_CORES)

	var similarities []Similarity
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Amount of computations that will be made
	wg.Add(len(mainStrings) * len(otherStrings))

	for _, mainString := range mainStrings {
		for _, otherString := range otherStrings {
			// Try to write to the channel, if it is full,
			// it will wait until it is free.
			waitChan <- struct{}{}

			// Now send a goroutine to calculate the similarity
			go func(mainString string, otherString string) {
				// Check if case insensitive
				if BoolFlags["Insensitive"] {
					mainString = strings.ToLower(mainString)
					otherString = strings.ToLower(otherString)
				}

				score := calculateSimilarity(mainString, otherString)
				similarity := Similarity{
					Metric: cases.Title(language.Und, cases.NoLower).String(Metric),
					S1:     mainString,
					S2:     otherString,
					Score:  score,
				}

				// Lock the array and append the similarity
				mu.Lock()
				similarities = append(similarities, similarity)
				mu.Unlock()
				wg.Done()

				// Unlock the channel
				<-waitChan
			}(mainString, otherString)
		}
	}

	wg.Wait()

	// Sort the slice by s1
	sort.Slice(similarities, func(i, j int) bool {
		return strings.ToLower(similarities[i].S1) < strings.ToLower(similarities[j].S1)
	})

	// Now check if its not set to silent to print results
	if !BoolFlags["Silent"] {
		for _, similarity := range similarities {
			similarity.Result()
		}
	}

	// Check if output to write to file
	if StringFlags["Output"] != "" {
		// Write to file
		writeToJson(StringFlags["Output"], similarities)
	}
}

// Flow for big files, for which I will not hold
// the similarities slice in memory and I'll
// be apending each result to the final json file.
func BigFileFlow(mainStrings []string, otherStrings []string, StringFlags map[string]string, BoolFlags map[string]bool) {
	// Getting the function based on the metric to use
	var Metric string
	var calculateSimilarity func(string, string) float64

	// Decide on the metric to use if has a flag
	if StringFlags["Metric"] != "" {
		// Make sure Metric is all lower case
		Metric = strings.ToLower(StringFlags["Metric"])
	} else {
		Metric = "jaro"
	}
	calculateSimilarity = GetSimilarityFunc(Metric)

	// Create a channel to guarantee max amount of goroutines
	// equal to cpu cores
	MAX_CPU_CORES := runtime.NumCPU()
	waitChan := make(chan struct{}, MAX_CPU_CORES)

	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(mainStrings) * len(otherStrings))

	// Create a json file with an empty array
	createEmptyJsonArrayFile(StringFlags["Output"])

	// Open the file
	file, err := os.OpenFile(StringFlags["Output"], os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Here, the file will start as [] and I'll be appending
	// each result to the file. I'll use a mutex to lock the file
	// while appending.
	// It's important to note that the first time I'm appending
	// I'll have to ommit a comma, since the normal apending is
	// `,{"key":value}]`.
	isEmpty := isEmptyList(file)

	for _, mainString := range mainStrings {
		for _, otherString := range otherStrings {
			// Try to write to the channel, if it is full,
			// it will wait until it is free.
			waitChan <- struct{}{}

			// Now send a goroutine to calculate the similarity
			go func(mainString string, otherString string) {
				// Check if case insensitive
				if BoolFlags["Insensitive"] {
					mainString = strings.ToLower(mainString)
					otherString = strings.ToLower(otherString)
				}

				score := calculateSimilarity(mainString, otherString)
				similarity := Similarity{
					Metric: cases.Title(language.Und, cases.NoLower).String(Metric),
					S1:     mainString,
					S2:     otherString,
					Score:  score,
				}

				// Lock the array and append the similarity
				mu.Lock()
				appendToJsonArray(file, &similarity, isEmpty)
				// After the first append, I'll set isEmpty to false
				if isEmpty {
					isEmpty = false
				}
				mu.Unlock()
				wg.Done()

				// Unlock the channel
				<-waitChan
			}(mainString, otherString)
		}
	}

	wg.Wait()

}
