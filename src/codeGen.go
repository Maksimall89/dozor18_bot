package src

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func CodeGen(prefix string, postfix string, maxCount int, pathFileWord string) string {

	if maxCount == 0 {
		return "&#10071;Слишком мало кодов: 0 кодов. \n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>"
	}

	if maxCount > 124 {
		maxCount = 124
	}

	rand.Seed(time.Now().UTC().UnixNano()) // real random

	var str string
	var code, buff, codePostfix, count int
	result := fmt.Sprintf("&#9989;Готовые коды (%d штук).\nКОД\tКО\tСектор\n\n", maxCount)
	arrCode := make(map[string]int)

	// read codes from file
	lines, err := readLines(pathFileWord)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	for count < maxCount {
		// if we have the same digital codes
		buff = code
		code = rand.Intn(90) + 10
		if code == buff {
			continue
		}

		// choose variant codes
		if postfix == "" && prefix == "" {
			buff = codePostfix
			codePostfix = rand.Intn(len(lines)-5) + 1
			if codePostfix == buff {
				continue
			}
			str = fmt.Sprintf("%s%d\t1\t1\r\n", lines[codePostfix], code)
		} else {
			codePostfix = rand.Intn(9) + 1
			str = fmt.Sprintf("%s%d%s%d\t1\t1\r\n", prefix, code, postfix, codePostfix)
		}

		// add new code
		if _, ok := arrCode[str]; !ok {
			arrCode[str] = count
			result += str
			count++
		}
	}
	return result
}
