package src

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

func TestReadLines(t *testing.T) {
	t.Parallel()

	// positive test
	lines, err := readLines(fmt.Sprintf("../%s", NameFileWords))
	if err != nil {
		t.Errorf("Can't open file %s", NameFileWords)
	}
	if len(lines) <= 3 {
		t.Errorf("File %s is empty", NameFileWords)
	}

	numberLine := rand.Intn(len(lines)-1) + 1
	if strings.ContainsAny(lines[numberLine], "") {
		t.Errorf(`Line "%s" in file %s is empty`, lines[numberLine], NameFileWords)
	}

	// negative test
	_, err = readLines("")
	if err == nil {
		t.Errorf("Can't open file %s", NameFileWords)
	}
}

func TestCodeGen(t *testing.T) {
	t.Parallel()

	type inputTest struct {
		prefix   string
		postfix  string
		maxCount int
		pathFile string
	}

	type testPair struct {
		input  inputTest
		output string
	}

	file := fmt.Sprintf("../%s", NameFileWords)

	var tests = []testPair{
		{inputTest{"", "", 0, file}, `Слишком мало кодов: 0 кодов. \n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>`},
		{inputTest{"", "", 1, file}, `&#9989;Готовые коды \(1 штук\)\.\nКОД\tКО\tСектор\n\n\W+\d+\t1\t1`},
		{inputTest{"", "", 1000, file}, `&#9989;Готовые коды \(124 штук\)\.\nКОД\tКО\tСектор\n\n\W+\d+\t1\t1`},
		{inputTest{"1D", "R", 0, file}, `Слишком мало кодов: 0 кодов. \n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>`},
		{inputTest{"1D", "R", 1, file}, `&#9989;Готовые коды \(1 штук\)\.\nКОД\tКО\tСектор\n\n1D\d+R\d+\t1\t1`},
		{inputTest{"1D", "R", 100, file}, `&#9989;Готовые коды \(100 штук\)\.\nКОД\tКО\tСектор\n\n1D\d+R\d+\t1\t1`},
		{inputTest{"2D", "", 1, file}, `&#9989;Готовые коды \(1 штук\)\.\nКОД\tКО\tСектор\n\n2D\d+\t1\t1`},
		{inputTest{"", "RR", 1, file}, `&#9989;Готовые коды \(1 штук\)\.\nКОД\tКО\tСектор\n\n\d+RR\d+\t1\t1`},
	}

	for _, pair := range tests {
		result := CodeGen(pair.input.prefix, pair.input.postfix, pair.input.maxCount, pair.input.pathFile)
		if match, _ := regexp.MatchString(pair.output, result); !match {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
