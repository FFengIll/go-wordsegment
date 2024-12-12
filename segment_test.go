package wordsegment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegmenter_Divide(t *testing.T) {
	input := "thisisatest"
	divided := segmenter.Divide(input)

	want := "[{t hisisatest} {th isisatest} {thi sisatest} {this isatest} {thisi satest} {thisis atest} {thisisa test} {thisisat est} {thisisate st} {thisisates t} {thisisatest }]"

	assert.Equal(t, want, fmt.Sprintf("%v", divided))
}
