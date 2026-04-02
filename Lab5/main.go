package main

import (
	"fmt"
	"strings"
)

func main() {
	var phrase string = "КОМПОЗИТ"
	var key string = "КОМБАЙН"

	res := encode(phrase, key, 4, 16)
	fmt.Println(string(res))
	res = decode(string(res), key, 4, 16)
	fmt.Println(string(res))
}

func encode(phrase string, key string, blockSize int, count int) []rune {
	phraseR := []rune(phrase)
	result := []rune(phrase)
	var left_part = []rune(strings.Repeat("0", blockSize))
	var right_part = []rune(strings.Repeat("0", blockSize))
	var temp []rune

	keyR := []rune(key)
	roundKeys := make([]rune, count)
	for i := 0; i < count; i++ {
		roundKeys[i] = keyR[i%len(keyR)]
	}
	for d := 0; d < len(phraseR)-len(phraseR)%(2*blockSize); d += 2 * blockSize {
		left_part = []rune(phraseR[d:min(d+blockSize, len(phraseR))])
		right_part = []rune(phraseR[min(d+blockSize, len(phraseR)):min(d+2*blockSize, len(phraseR))])
		//fmt.Println("before round", string(left_part), string(right_part))

		for i := 0; i < count; i++ {
			temp = xor(left_part, F(right_part, roundKeys[i]))
			left_part = right_part
			right_part = temp
			//fmt.Println("round", i, string(left_part), string(right_part))
		}
		result = insert(result, right_part, d)
		result = insert(result, left_part, d+blockSize)
	}
	return result
}

func xor(subblock []rune, key []rune) []rune {
	result := make([]rune, len(subblock))
	keyR := []rune(key)
	for i := 0; i < len(subblock); i++ {
		result[i] = subblock[i] ^ keyR[i%len(keyR)]
	}
	return result
}

func F(subblock []rune, key rune) []rune {
	result := make([]rune, len(subblock))
	for i := range len(result) {
		result[i] = subblock[i] ^ key
	}
	return result
}

func insert(dest []rune, val []rune, startIndex int) []rune {
	result := []rune(dest)
	for i := 0; i < min(len(val), (len(dest)-startIndex)); i++ {
		result[startIndex+i] = val[i]
	}
	return result
}

func decode(phrase string, key string, blockSize int, count int) []rune {
	phraseR := []rune(phrase)
	result := []rune(phrase)
	var left_part = []rune(strings.Repeat("0", blockSize))
	var right_part = []rune(strings.Repeat("0", blockSize))
	var temp []rune

	keyR := []rune(key)
	roundKeys := make([]rune, count)
	for i := 0; i < count; i++ {
		roundKeys[i] = keyR[i%len(keyR)]
	}
	for d := 0; d < len(phraseR)-len(phraseR)%(2*blockSize); d += 2 * blockSize {

		left_part = []rune(phraseR[d : d+blockSize])
		right_part = []rune(phraseR[d+blockSize : d+2*blockSize])
		//fmt.Println("before round", string(left_part), string(right_part))

		for i := count - 1; i >= 0; i-- {
			temp = xor(left_part, F(right_part, roundKeys[i]))
			left_part = right_part
			right_part = temp
			//fmt.Println("round", i, string(left_part), string(right_part))
		}
		result = insert(result, right_part, d)
		result = insert(result, left_part, d+blockSize)
	}
	return result
}
