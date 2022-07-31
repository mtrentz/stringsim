package similarity

import (
	"fmt"

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
