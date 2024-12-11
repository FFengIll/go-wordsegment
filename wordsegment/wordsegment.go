package wordsegment

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"math"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed unigrams.txt bigrams.txt words.txt
var content embed.FS

// Segmenter represents the word segmentation logic
type Segmenter struct {
	Unigrams map[string]float64
	Bigrams  map[string]float64
	Total    float64
	Limit    int
	Words    []string
	Alphabet set
}

// set is a helper type to represent a set of characters
type set map[rune]struct{}

// NewSegmenter initializes a new Segmenter
func NewSegmenter() *Segmenter {
	return &Segmenter{
		Unigrams: make(map[string]float64),
		Bigrams:  make(map[string]float64),
		Limit:    24,
		Alphabet: set{'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {}, 'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {}, 'm': {}, 'n': {}, 'o': {}, 'p': {}, 'q': {}, 'r': {}, 's': {}, 't': {}, 'u': {}, 'v': {}, 'w': {}, 'x': {}, 'y': {}, 'z': {}, '0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}},
	}
}

// Load loads the unigram, bigram, and word data into the Segmenter
func (s *Segmenter) Load() error {
	unigramsData, err := content.ReadFile("unigrams.txt")
	if err != nil {
		logrus.Fatalf("Error reading unigrams.txt: %v", err)
	}
	s.Unigrams = parse(unigramsData)

	bigramsData, err := content.ReadFile("bigrams.txt")
	if err != nil {
		logrus.Fatalf("Error reading bigrams.txt: %v", err)
	}
	s.Bigrams = parse(bigramsData)

	wordsData, err := content.ReadFile("words.txt")
	if err != nil {
		logrus.Fatalf("Error reading words.txt: %v", err)
	}
	s.Words = strings.Split(string(wordsData), "\n")

	s.Total = 1024908267229.0
	return nil
}

// parse parses a byte slice into a map of words to counts
func parse(data []byte) map[string]float64 {
	result := make(map[string]float64)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) == 2 {
			var word string
			var count float64
			fmt.Sscanf(parts[0], "%s", &word)
			fmt.Sscanf(parts[1], "%f", &count)
			result[word] = count
		}
	}
	if err := scanner.Err(); err != nil {
		logrus.Fatalf("Error parsing file: %v", err)
	}
	return result
}

// Score computes the probability score of a word given the previous word
func (s *Segmenter) Score(word, previous string) float64 {
	if previous == "" {
		if score, found := s.Unigrams[word]; found {
			return score / s.Total
		}
		return 10.0 / (s.Total * math.Pow(10, float64(len(word))))
	}

	bigram := fmt.Sprintf("%s %s", previous, word)
	if score, found := s.Bigrams[bigram]; found {
		if prevScore, found := s.Unigrams[previous]; found {
			return score / s.Total / prevScore
		}
	}
	return s.Score(word, "")
}

// Isegment returns the best segmentation of the text using a dynamic programming approach
func (s *Segmenter) Isegment(text string) []string {
	memo := make(map[string][]string)

	var search func(text, previous string) (float64, []string)

	search = func(text, previous string) (float64, []string) {
		if text == "" {
			return 0.0, nil
		}

		var candidates []struct {
			score float64
			words []string
		}

		for _, prefixSuffix := range s.Divide(text) {
			prefix, suffix := prefixSuffix[0], prefixSuffix[1]
			prefixScore := math.Log10(s.Score(prefix, previous))
			var suffixScore float64
			var suffixWords []string
			if cached, found := memo[prefix+suffix]; found {
				suffixScore, suffixWords = 0.0, cached
			} else {
				suffixScore, suffixWords = search(suffix, prefix)
			}

			candidates = append(candidates, struct {
				score float64
				words []string
			}{prefixScore + suffixScore, append([]string{prefix}, suffixWords...)})
		}

		if len(candidates) <= 0 {
			return 0.0, nil
		}

		// println(candidates)
		best := candidates[0]
		for _, candidate := range candidates {
			if candidate.score > best.score {
				best = candidate
			}
		}

		memo[text+previous] = best.words
		return best.score, best.words
	}

	// Clean the text (equivalent to Python's clean())
	cleanedText := s.Clean(text)

	// Define chunk size and initialize prefix
	size := 250
	var prefix string
	var result []string

	// Loop through the cleaned text in chunks of `size`
	for offset := 0; offset < len(cleanedText); offset += size {
		// Get the current chunk
		end := offset + size
		if end > len(cleanedText) {
			end = len(cleanedText)
		}
		chunk := cleanedText[offset:end]

		// Combine prefix and chunk, then call search to segment
		combinedText := prefix + chunk
		_, chunkWords := search(combinedText, "<s>")

		// Update the prefix with the last 5 words from the chunk
		if len(chunkWords) > 5 {
			prefix = strings.Join(chunkWords[len(chunkWords)-5:], "")
			// Add all the words from the chunk (except the last 5 words which are used as prefix)
			for _, word := range chunkWords[:len(chunkWords)-5] {
				result = append(result, word)
			}
		} else {
			prefix = strings.Join(chunkWords, "")
			// Add all the words from the chunk (except the last 5 words which are used as prefix)
			for _, word := range chunkWords {
				result = append(result, word)
			}
		}
	}

	// Finally, process any remaining prefix words (from the last chunk)
	_, prefixWords := search(prefix, "<s>")
	for _, word := range prefixWords {
		result = append(result, word)
	}
	return result
}

// Divide splits the text into all possible prefix-suffix pairs up to the limit
func (s *Segmenter) Divide(text string) [][2]string {
	var result [][2]string
	for pos := 1; pos < len(text) && pos <= s.Limit; pos++ {
		result = append(result, [2]string{text[:pos], text[pos:]})
	}
	return result
}

// Clean cleans up the input text by removing non-alphanumeric characters
func (s *Segmenter) Clean(text string) string {
	var result strings.Builder
	for _, r := range text {
		if _, found := s.Alphabet[r]; found {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// func main() {
// 	segmenter := NewSegmenter()
// 	if err := segmenter.Load(); err != nil {
// 		logrus.Fatalf("Error loading segmenter: %v", err)
// 	}

// 	input := "thisisatest"
// 	segmented := segmenter.Isegment(input)

// 	for _, word := range segmented {
// 		logrus.Infof("Segmented word: %s", word)
// 	}
// }
