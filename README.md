# bGOw

A small, dependency-free **Bag-of-Words** implementation in Go.

`bGOw` turns text into fixed-length binary vectors. You fit an `Engine` on a
corpus to build a vocabulary, then transform any sentence into a `[]uint8`
vector where each position indicates whether the corresponding vocabulary word
is present. It can also map vectors back to the words they represent.

## Why

I needed a simple bag-of-words implementation in Go for my startup. :)

## Install

```bash
go get github.com/mriceman/bgow@latest
```

```go
import "github.com/mriceman/bgow/pkg/bow"
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/mriceman/bgow/pkg/bow"
)

func main() {
	corpus := []string{
		"my name is martin",
		"i have a dog",
		"i like coding",
	}

	engine := bow.NewEngine().WithCorpus(corpus)
	engine.Fit()

	// Encode a sentence into a binary bag-of-words vector.
	vec := engine.Transform("martin is NOT a dog")
	fmt.Println(vec) // e.g. [1 0 1 1 0 ...]

	// Decode the vector back into the words it represents.
	fmt.Println(engine.GetWordsFromVec(vec)) // [martin is a dog]

	// Look up the position of a word in the vocabulary (-1 if absent).
	fmt.Println(engine.GetPos("dog"))
}
```

## How it works

1. **Fit** splits every entry in the corpus on spaces and counts word
   frequencies, then assigns each word a fixed position in the vocabulary. It
   also precomputes the binary vector for every corpus entry (`CorpusVecs`) and
   the per-word frequency counts (`CountVec`).
2. **Transform** produces a vector the size of the vocabulary, setting `1` at
   the position of every known word found in the input (unknown words are
   ignored). The vector is binary (presence/absence), not a frequency count.

## API

| Method | Description |
| --- | --- |
| `NewEngine() *Engine` | Creates a new engine. |
| `WithCorpus(corpus []string) *Engine` | Sets the corpus (chainable). |
| `WithTop(top int) *Engine` | Limits the vocabulary to the `top` most frequent words (chainable). `0` means use all words. |
| `SetCorpus(corpus []string)` | Sets the corpus without chaining. |
| `Fit()` | Builds the vocabulary from the corpus and precomputes corpus vectors and word counts. |
| `Transform(payload string) []uint8` | Encodes a sentence into a binary vector. |
| `BowLength() *Length` | Returns the vector length (`Length.X`). |
| `GetPos(word string) int` | Returns a word's vocabulary index, or `-1` if not present. |
| `LookUpWord(idx int) (string, error)` | Returns the word at a vocabulary index, or an error if out of range. |
| `GetWordsFromVec(vec []uint8) []string` | Decodes a vector back into words. |
| `CorpusVecs() [][]uint8` | Returns the precomputed binary vectors for each corpus entry (after `Fit`). |
| `CountVec() []int` | Returns the frequency count for each vocabulary word, ordered by position (after `Fit`). |

### Limiting vocabulary size

Use `WithTop` to keep only the most frequent words, which caps the vector size:

```go
engine := bow.NewEngine().WithCorpus(corpus).WithTop(1000)
engine.Fit()
```

## Testing

```bash
go test ./...
```
```
