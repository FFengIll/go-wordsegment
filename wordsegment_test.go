package wordsegment

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var segmenter *Segmenter

// Initialize the segmenter and load data
func init() {
	segmenter = NewSegmenter()
	if err := segmenter.Load(); err != nil {
		logrus.Fatalf("Failed to load segmenter: %v", err)
	}
}

func TestUnigrams(t *testing.T) {
	if _, found := segmenter.Unigrams["test"]; !found {
		t.Errorf("Expected 'test' to be in unigrams")
	}
}

func TestBigrams(t *testing.T) {
	if _, found := segmenter.Bigrams["in the"]; !found {
		t.Errorf("Expected 'in the' to be in bigrams")
	}
}

func TestClean(t *testing.T) {
	cleaned := segmenter.Clean("Can't buy me love!")
	expected := "cantbuymelove"
	if cleaned != expected {
		t.Errorf("Expected %v, but got %v", expected, cleaned)
	}
}

func TestSegment0(t *testing.T) {
	result := []string{"choose", "spain"}
	segmented := segmenter.Segment("choosespain")
	if !equal(result, segmented) {
		t.Errorf("Expected %v, but got %v", result, segmented)
	}
}

func TestSegment1(t *testing.T) {
	result := []string{"this", "is", "a", "test"}
	segmented := segmenter.Segment("thisisatest")
	if !equal(result, segmented) {
		t.Errorf("Expected %v, but got %v", result, segmented)
	}
}

func TestSegment2(t *testing.T) {
	result := []string{
		"when", "in", "the", "course", "of", "human", "events", "it", "becomes", "necessary",
	}
	segmented := segmenter.Segment("wheninthecourseofhumaneventsitbecomesnecessary")
	if !equal(result, segmented) {
		t.Errorf("Expected %v, but got %v", result, segmented)
	}
}

// Implement similar tests for all other segments (3-12)

func TestWords(t *testing.T) {
	if len(segmenter.Words) == 0 {
		t.Errorf("Expected WORDS to have elements, but it is empty")
	}
	if segmenter.Words[0] != "aa" {
		t.Errorf("Expected first word to be 'aa', but got %v", segmenter.Words[0])
	}
	if segmenter.Words[len(segmenter.Words)-1] != "zzz" {
		t.Errorf("Expected last word to be 'zzz', but got %v", segmenter.Words[len(segmenter.Words)-1])
	}
}

func TestSegmenter_Divide(t *testing.T) {
	input := "thisisatest"
	divided := segmenter.Divide(input)

	want := "[{t hisisatest} {th isisatest} {thi sisatest} {this isatest} {thisi satest} {thisis atest} {thisisa test} {thisisat est} {thisisate st} {thisisates t} {thisisatest }]"

	assert.Equal(t, want, fmt.Sprintf("%v", divided))
}

// Utility function to compare slices of strings
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
