package main

import (
	. "github.com/FFengIll/go-wordsegment/wordsegment"
	"github.com/sirupsen/logrus"
)

func main() {
	segmenter := NewSegmenter()
	if err := segmenter.Load(); err != nil {
		logrus.Fatalf("Error loading segmenter: %v", err)
	}

	input := "thisisatest"
	segmented := segmenter.Isegment(input)

	for _, word := range segmented {
		logrus.Infof("Segmented word: %s", word)
	}
}
