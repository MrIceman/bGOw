package bow_test

import (
	"bGOw/pkg/bow"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	t.Run("should set BowLength properly", func(t *testing.T) {
		n := 10000
		var corpus []string
		for i := 0; i < n; i++ {
			corpus = append(corpus, fmt.Sprintf("%d,", i))
		}

		subject := bow.NewEngine().WithCorpus(corpus)

		subject.Fit()

		shape := subject.BowLength()

		if shape.X != n {
			t.Errorf("shape.X != n, got %v", shape.X)
		}
	})

	t.Run("should return proper vector", func(t *testing.T) {
		corpus := []string{
			"my name is martin",
			"i have a dog",
			"i like coding",
		}
		subject := bow.NewEngine().WithCorpus(corpus)

		subject.Fit()

		length := subject.BowLength()
		if length.X != 10 {
			t.Errorf("length.X != 10, got %v", length.X)
		}

		payload := "martin is NOT a dog"
		result := subject.Transform(payload)

		var tokens []string
		for _, token := range subject.GetWordsFromVec(result) {
			slog.Info("reverted token", "token", token, "pos", subject.GetPos(token))
			tokens = append(tokens, token)
		}
		if len(tokens) != 4 {
			t.Errorf("len(tokens)=%d, len(tokens)=%d", len(tokens), len(tokens))
		}

		for _, token := range tokens {
			if !strings.Contains(payload, token) {
				t.Errorf("token %s not found in payload", token)
			}
		}
	})

	t.Run("should convert a big ass vector", func(t *testing.T) {
		y := 380000
		x := 3000
		var sentences []string
		for i := 0; i < x; i++ {
			sentences = append(sentences, fmt.Sprintf("%d,", i))
		}
		var corpus []string
		log.Println("preparing corpus")
		for i := 0; i < y; i++ { // corpus has y entries
			for j := 0; j < x; j++ { // each "sentence" has x tokens
				corpus = append(corpus, fmt.Sprintf("%d,", j))
				log.Printf("i: %d, j: %d", i, j)
			}
		}

		log.Println("starting fit")
		v := time.Now()
		subject := bow.NewEngine().WithCorpus(corpus)
		subject.Fit()
		w := time.Now()
		slog.Info("fit done", "duration", w.Sub(v).Seconds())
	})
}
