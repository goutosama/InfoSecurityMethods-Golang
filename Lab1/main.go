package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func keywordOrder(keyword string) []int {
	runes := []rune(keyword)

	type pair struct {
		char  rune
		index int
	}

	pairs := make([]pair, len(runes))

	for i, r := range runes {
		pairs[i] = pair{char: r, index: i}
	}

	// сортируем по символу
	sort.SliceStable(pairs, func(i, j int) bool {
		if pairs[i].char == pairs[j].char { // если одинаковые - по порядку индекса
			return pairs[i].index < pairs[j].index
		}
		return pairs[i].char < pairs[j].char
	})

	result := make([]int, len(runes))

	for i, p := range pairs {
		result[i] = p.index
	}

	return result
}

func main() {
	// обработка флагов командной строки
	// -enc для кодирования, -dec - декодирования
	// -codeword кодовое слово
	// последний параметр - шифруемая фраза
	isEncode := flag.Bool("enc", false, "Encode message with a codeword specified with --codeword")
	isDecode := flag.Bool("dec", false, "Decode message with a codeword specified with --codeword")
	codeword := flag.String("codeword", "", "A codeword to use with a phrase")

	flag.Parse()

	if len(os.Args) > 1 {
		var phrase string = os.Args[len(os.Args)-1]

		if *codeword != "" {
			if !(*isEncode && *isDecode) {
				if *isEncode {
					fmt.Println(encode([]rune(*codeword), []rune(phrase)))
				}
				if *isDecode {
					fmt.Println(decode([]rune(*codeword), []rune(phrase)))
				}
			}
		}
	}
}

func encode(codeword []rune, phrase []rune) string {
	var codewordLen int = len(codeword)
	var fullCol int = len(phrase) / codewordLen
	var partColLen int = len(phrase) % len(codeword)

	var order []int = keywordOrder(string(codeword))

	result := make([]rune, len(phrase))
	copy(result, phrase)

	var index int = 0
	for i := range codewordLen {
		var columnSize int = fullCol
		if order[i] < partColLen {
			columnSize = fullCol + 1
		}
		for j := range columnSize {
			result[index] = phrase[order[i]+(codewordLen)*j]
			index++
		}
	}
	return string(result)
}

func decode(codeword []rune, phrase []rune) string {
	var codewordLen int = len(codeword)
	var fullCol int = len(phrase) / codewordLen
	var partColLen int = len(phrase) % len(codeword)

	var order []int = keywordOrder(string(codeword))

	result := make([]rune, len(phrase))
	copy(result, phrase)

	var index int = 0
	for i := range codewordLen {
		var columnSize int = fullCol
		if order[i] < partColLen {
			columnSize++
		}
		for j := range columnSize {
			result[order[i]+codewordLen*j] = phrase[index]
			index++
		}
	}

	return string(result)
}
