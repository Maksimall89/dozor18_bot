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
		return "Слишком мало кодов: 0 кодов. \n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>"
	}

	if maxCount > 124 {
		maxCount = 124
	}

	rand.Seed(time.Now().UTC().UnixNano()) // real random

	var str string
	var code, buff, codePostfix, count int

	arrCode := make(map[string]string)

	// read codes from file
	lines, err := readLines(pathFileWord)
	if err != nil {
		log.Println(err)
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
			str = lines[codePostfix]
			codePostfix = rand.Intn(len(lines)-5) + 1
			if lines[codePostfix] == str {
				continue
			}
			str = fmt.Sprintf("%s%d\t1\t1\r\n", lines[codePostfix], code)
		} else {
			codePostfix = rand.Intn(9) + 1
			str = fmt.Sprintf("%s%d%s%d\t1\t1\r\n", prefix, code, postfix, codePostfix)
		}

		// add new code
		if _, ok := arrCode[str]; !ok {
			arrCode[str] = str
			count++
		}
	}

	str = fmt.Sprintf("&#9989;Готовые коды (%d штук).\nКОД\tКО\tСектор\n\n", maxCount)
	for key := range arrCode {
		str += key
	}

	return str
}
