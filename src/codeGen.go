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

func CodeGen(prefix string, postfix string, typeStyle bool, maxCount int) string {

	if maxCount == 0 {
		return "Слишком мало кодов: 0 кодов. \n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>"
	}

	if maxCount > 124 {
		maxCount = 124
	}

	rand.Seed(time.Now().UTC().UnixNano()) // real random

	var str string
	var code, buff, codePostfix, count int
	var isSame bool

	arrCode := make(map[int]string)

	// read codes from file
	lines, err := readLines(NameFileWords)
	if err != nil {
		log.Println(err)
	}

	for {
		// if we have the same digital codes
		buff = code
		code = rand.Intn(90) + 10
		if code == buff {
			continue
		}

		// choose variant codes
		if typeStyle {
			codePostfix = rand.Intn(9) + 1
			str = fmt.Sprintf("%s%d%s%d\t1\t1\r\n", prefix, code, postfix, codePostfix)
		} else {
			codePostfix = rand.Intn(len(lines)-5) + 1
			str = fmt.Sprintf("%s%d\t1\t1\r\n", lines[codePostfix], code)
		}

		// check old codes
		isSame = false

		for _, value := range arrCode {
			if value == str {
				isSame = true
				break
			}
		}

		// we not have double code
		if isSame {
			continue
		}

		// add new code
		arrCode[count] = str

		// check count codes
		if len(arrCode) == maxCount {
			break
		} else {
			count++
		}
	}

	str = fmt.Sprintf("&#9989;Готовые коды (%d штук).\nКОД\tКО\tСектор\n\n", maxCount)
	for _, value := range arrCode {
		str += value
	}

	return str
}
