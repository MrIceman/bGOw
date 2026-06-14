package bow

import (
	"sort"
	"strings"
)

type (
	Length struct {
		X int
	}
	Engine struct {
		corpus     []string
		countVec   []int
		corpusVecs [][]uint8
		// top specifies the top most occurrences of a word within the corpus that should be considered in the bow
		// otherwise ALL words will be taken
		top        int
		vectorSize int

		dict       map[string]int
		posLookUp  map[int]string // has the position of the word in the dict
		wordLookUp map[string]int
		shape      *Length
	}
)

func NewEngine() *Engine {
	return &Engine{
		dict:       make(map[string]int),
		posLookUp:  make(map[int]string),
		wordLookUp: make(map[string]int),
	}
}

func (e *Engine) BowLength() *Length {
	return e.shape
}

func (e *Engine) WithCorpus(corpus []string) *Engine {
	e.corpus = corpus
	return e
}

func (e *Engine) WithTop(top int) *Engine {
	e.top = top
	return e
}

func (e *Engine) SetCorpus(corpus []string) {
	e.corpus = corpus
}

type count = int
type wordCount struct {
	Word  string
	Count int64
}

// GetPos returns the position of the word in the dictionary, if not exists then -1
func (e *Engine) GetPos(word string) int {
	pos, ok := e.wordLookUp[word]
	if !ok {
		return -1
	}
	return pos
}

func (e *Engine) GetWordsFromVec(vec []uint8) []string {
	var result []string
	for idx, val := range vec {
		if val == 1 {
			result = append(result, e.posLookUp[idx])
		}
	}

	return result
}

func (e *Engine) Fit() {
	countWMap := make(map[string]int64)
	for cIdx := range e.corpus {
		for _, w := range strings.Split(e.corpus[cIdx], " ") {
			countWMap[w]++
		}
	}

	// initialize vector length -- it's either top or the vec size
	var vecSize int
	if len(countWMap) < e.top || e.top == 0 {
		vecSize = len(countWMap)
	} else {
		vecSize = e.top
	}
	e.vectorSize = vecSize

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
		e.posLookUp[i] = wordCounts[i].Word
		e.wordLookUp[wordCounts[i].Word] = i
		e.countVec[i] = int(wordCounts[i].Count)
	}

	e.shape = &Length{
		X: e.vectorSize,
	}
	// create the corpus vector

	for _, c := range e.corpus {
		vec := e.Transform(c)
		e.corpusVecs = append(e.corpusVecs, vec)
	}
}

func (e *Engine) CorpusVecs() [][]uint8 {
	return e.corpusVecs
}

func (e *Engine) CountVec() []int {
	return e.countVec
}

func (e *Engine) Transform(payload string) []uint8 {
	initVector := make([]uint8, e.vectorSize)
	splitPayload := strings.Split(payload, " ")
	for idx := range splitPayload {
		pos, ok := e.wordLookUp[splitPayload[idx]]
		if !ok {
			continue
		}
		initVector[pos] = 1
	}

	return initVector
}
