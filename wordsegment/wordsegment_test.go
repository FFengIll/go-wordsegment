package wordsegment

import (
	"testing"

	"github.com/sirupsen/logrus"
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

// func TestMain(t *testing.T) {
// 	// Redirect stdout to capture output
// 	old := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w

// 	// Create a test file
// 	err := ioutil.WriteFile("test_input.txt", []byte("choose spain\nthis is a test\n"), 0644)
// 	if err != nil {
// 		t.Fatalf("Error writing test input file: %v", err)
// 	}
// 	defer os.Remove("test_input.txt")

// 	// Call the main function with the file
// 	main([]string{"test_input.txt"})

// 	// Close the writer and read the output
// 	w.Close()
// 	var buf bytes.Buffer
// 	_, err = buf.ReadFrom(r)
// 	if err != nil {
// 		t.Fatalf("Error reading from pipe: %v", err)
// 	}

// 	// Compare the result
// 	expected := "choose spain\nthis is a test\n"
// 	if buf.String() != expected {
// 		t.Errorf("Expected %v, but got %v", expected, buf.String())
// 	}

// 	// Restore stdout
// 	os.Stdout = old
// }

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
