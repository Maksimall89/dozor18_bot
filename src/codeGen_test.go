package src

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestReadLines(t *testing.T) {
	t.Parallel()

	// positive test
	_, err := readLines(fmt.Sprintf("../%s", NameFileWords))
	if err != nil {
		t.Errorf("Can't open file %s", NameFileWords)
	}
	// negative test
	_, err = readLines("")
	if err == nil {
		t.Errorf("Can't open file %s", NameFileWords)
	}

}

func TestCodeGen(t *testing.T) {
	t.Parallel()

	lines, _ := readLines(fmt.Sprintf("../%s", NameFileWords))
	if len(lines) <= 3 {
		t.Errorf("File %s is empty", NameFileWords)
	}
	numberLine := rand.Intn(len(lines)-1) + 1
	if strings.ContainsAny(lines[numberLine], "") {
		t.Errorf(`Line "%s" in file %s is empty`, lines[numberLine], NameFileWords)
	}
}
