package commons

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func ReadWords(p []byte) (words []string) {
	words = make([]string, 0, 1)
	adv := 0
	for {
		adv0, tkn, err := bufio.ScanWords(p[adv:], true)
		if err != nil {
			fmt.Println("read words from line failed,", string(p))
			os.Exit(9)
		}
		if tkn == nil {
			break
		}
		adv = adv + adv0
		words = append(words, string(tkn))
	}
	return
}

func WordsToLine(words []string) (line string) {
	buf := bytes.NewBufferString(line)
	for _, word := range words {
		buf.WriteString(word)
		buf.WriteString(" ")
	}
	line = strings.TrimSpace(buf.String())
	return
}

func WordsIndex(words []string, word string) (idx int) {
	idx = -1
	word = strings.ToUpper(word)
	for i, word0 := range words {
		if strings.ToUpper(word0) == word {
			idx = i
			break
		}
	}
	return
}

func WordsContainsOne(words []string, sub ...string) (has bool) {
	for _, word := range sub {
		word = strings.ToUpper(word)
		for _, word0 := range words {
			if strings.ToUpper(word0) == word {
				has = true
				break
			}
		}
	}
	return
}

func WordsContainsAll(words []string, sub ...string) (has bool) {
	n := len(sub)
	if n == 0 {
		return
	}
	for _, word := range sub {
		word = strings.ToUpper(word)
		for _, word0 := range words {
			if strings.ToUpper(word0) == word {
				n--
			}
		}
	}
	has = n == 0
	return
}
