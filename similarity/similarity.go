package similarity

import (
	"fmt"
	"os"

	"github.com/antzucaro/matchr"
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

func GetSimilarity(s1 string, s2 string) Similarity {
	return Similarity{
		Metric: "Jaro",
		S1:     s1,
		S2:     s2,
		Score:  matchr.Jaro(s1, s2),
	}
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
