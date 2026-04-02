package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// обработка флагов командной строки
	// -codeword кодовое слово
	// последний параметр - шифруемая фраза
	codeword := flag.String("codeword", "", "A codeword to use with a phrase \n Encoding and decoding is the same process")

	flag.Parse()

	if len(os.Args) > 1 {
		var phrase string = os.Args[len(os.Args)-1]
		if *codeword != "" {
			fmt.Println(code(*codeword, phrase))
		}
	}
}

func code(codeword string, phrase string) string {
	var result = []byte(phrase)
	for i := range len(result) {
		result[i] = result[i] ^ codeword[i%len(codeword)]
	}
	return string(result)
}
