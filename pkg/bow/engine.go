package bow

import (
	"sort"
	"strings"
)

type Engine struct {
	Corpus []string
	Vector []uint8
	X      int
	Y      int
	Top    int

	dict       map[string]int
	posLookUp  map[int]string // has the position of the word in the dict
	wordLookUp map[string]int
}

func NewEngine() *Engine {
	return &Engine{
		dict:       make(map[string]int),
		posLookUp:  make(map[int]string),
		wordLookUp: make(map[string]int),
	}
}

func (e *Engine) SetCorpus(corpus []string) {
	e.Corpus = corpus
}

type count = int
type wordCount struct {
	Word  string
	Count int64
}

// GetPos returns the position of the word in the dictionary, if not exists then -1
func (s *Engine) GetPos(word string) int {
	pos, ok := s.wordLookUp[word]
	if !ok {
		return -1
	}
	return pos
}

func (s *Engine) GetWordsFromVec(vec []uint8) []string {
	var result []string
	for idx, val := range vec {
		if val == 1 {
			result = append(result, s.posLookUp[idx])
		}
	}

	return result
}

func (s *Engine) Fit() {
	countWMap := make(map[string]int64)
	for cIdx := range s.Corpus {
		for _, w := range strings.Split(s.Corpus[cIdx], " ") {
			countWMap[w]++
		}
	}
	var vecSize int
	if len(countWMap) < s.Top {
		vecSize = len(countWMap)
	} else {
		vecSize = s.Top
	}
	s.Vector = make([]uint8, vecSize)

	wordCounts := make([]wordCount, 0, len(countWMap))

	for word, count := range countWMap {
		wordCounts = append(wordCounts, wordCount{
			Word:  word,
			Count: count,
		})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].Count > wordCounts[j].Count
	})

	sortedWords := make(map[string]int, vecSize)

	for i := 0; i < vecSize; i++ {
		sortedWords[wordCounts[i].Word] = i
		s.posLookUp[i] = wordCounts[i].Word
		s.wordLookUp[wordCounts[i].Word] = i
	}
}

func (s *Engine) Transform(payload string) []uint8 {
	initVector := make([]uint8, len(s.Vector), len(s.Vector))
	splitPayload := strings.Split(payload, " ")
	for idx := range splitPayload {
		pos, ok := s.wordLookUp[splitPayload[idx]]
		if !ok {
			continue
		}
		initVector[pos] = 1
	}

	return initVector
}
