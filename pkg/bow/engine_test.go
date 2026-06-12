package bow

import (
	"log/slog"
	"testing"
)

func TestEngine(t *testing.T) {
	eng := NewEngine()
	corpus := []string{
		"This is a car",
		"This is a fish",
		"Car is broken at least",
		"what do you wish car",
		"car car car",
		"how u doing",
		"body",
	}
	eng.Top = 50000
	eng.SetCorpus(corpus)
	eng.Fit()

	vec := eng.Transform("The car is broken but at least I have a fish body")

	words := eng.GetWordsFromVec(vec)

	slog.Info("got a result", "res", vec, "vec count", len(vec), "words", words)

	for _, word := range words {
		slog.Info("look up", "word", word, "pos", eng.GetPos(word))
	}
}
