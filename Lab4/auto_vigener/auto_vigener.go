package auto_vigener

import (
	"flag"
	"fmt"
	"math"
	"os"
)

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
	var result []rune = []rune(phrase)
	var codeIndex int = 0
	for i := range phrase {
		if codeIndex < len(codeword) {
			result[i] = shiftLetter(phrase[i], codeword[i])
			codeIndex++
		} else {
			result[i] = shiftLetter(phrase[i], phrase[i])
		}
	}

	return string(result)
}

func decode(codeword []rune, phrase []rune) string {
	var result []rune = []rune(phrase)
	var codeIndex int = 0
	for i := range phrase {
		if codeIndex < len(codeword) {
			result[i] = shiftLetterBack(phrase[i], codeword[i])
			codeIndex++
		} else {
			result[i] = shiftLetterBack(phrase[i], phrase[i+codeIndex-1])
		}
	}
	return string(result)
}

func shiftLetter(letter rune, shift rune) rune {
	return (letter + shift) % int32(math.Pow(2, 16))
}

func shiftLetterBack(letter rune, shift rune) rune {
	var pow = int32(math.Pow(2, 16))
	return (pow + letter - shift) % pow
}
